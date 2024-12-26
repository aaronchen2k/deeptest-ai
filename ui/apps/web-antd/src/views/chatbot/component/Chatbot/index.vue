<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';

import { addSepIfNeeded } from '@vben/utils';

import { fetchEventSource } from '@microsoft/fetch-event-source';
import { message } from 'ant-design-vue';
import { Markdown } from 'vue3-markdown-it';

import consts from '#/config/constant';
import { getCache, setCache } from '#/utils/cache-local';

import {
  getDocLink,
  getLatestRobotMsg,
  isUnderRobotMsg,
  replaceLinkWithoutTitle,
  scroll,
  setSelectionRange,
} from './service';

const props = defineProps({
  llm: {
    type: String,
    default: '', // will use the first llm in chatchat server if empty
    required: false,
  },
  defaultKb: {
    type: String,
    default: 'wiki',
    required: true,
  },
  serverUrl: {
    type: String,
    default: 'http://127.0.0.1:9085/api/v1',
    required: true,
  },
});

const wikiAddress = 'https://wiki.deeptestcloud.com';
const wakeUpWord = '小深';
const humanName = 'Albert';
const humanAvatar = '/static/chat-einstein.png';

const CHAT_HISTORIES = 'chat_history_key';
const histories = ref([] as any[]);
const historyIndex = ref(-1);

const kb = ref(props.defaultKb);
const msg = ref('');
const isChatting = ref(false);
const continueOnCurrMsg = ref(false);

const mdText = `提取器用于将请求响应结果中的数据经过解析后提取出来，并存储在变量中，用于结果的校验或传递给下个请求调用等进一步操作。使用方法如下：1.响应头参数值提取：直接从响应头中提取所需参数值。2.响应体内容提取：-使用Xpath对响应体进行解析。-本系统使用jsonquery工具进行解析。-Xpath语法详见：https://github.com/antchfx/jsonquery示例XPath写法：\`\`\`xml<response-body><example><data>所需数据</data></example></response-body>\`\`\`XPath表达式可能如下：\`\`\`xpath/response-body/example/data\`\`\`请注意，具体的XPath表达式需要根据实际的响应体结构进行调整。
dom.ts:10 606`;
window.console.log(mdText);

const messages = ref([] as any[]);
messages.value.push(
  {
    type: 'human',
    name: humanName,
    content: wakeUpWord,
    avatar: humanAvatar,
  },
  {
    type: 'robot',
    name: 'ChatGPT',
    content: '你好，有什么可以帮助您的？',
    docs: '',
  },
);

const send = async () => {
  window.console.log('send ...');
  msg.value = msg.value.trim();
  if (!msg.value) return;

  const index = histories.value.indexOf(msg.value);
  if (index !== -1) {
    histories.value.splice(index, 1);
  }

  if (histories.value.length >= 30)
    histories.value = histories.value.splice(0, 1);

  const userMsg = msg.value;
  if (`${userMsg}` !== wakeUpWord) {
    histories.value.push(`${userMsg}`);
    historyIndex.value = histories.value.length;
    setCache(CHAT_HISTORIES, histories.value);
    msg.value = '';
  }

  isChatting.value = true;
  continueOnCurrMsg.value = false;

  const humanMsg = {
    type: 'human',
    name: humanName,
    content: userMsg,
    avatar: humanAvatar,
  };
  messages.value.push(humanMsg);
  scroll();

  const serverUrl = addSepIfNeeded(props.serverUrl);
  const url = `${serverUrl}aichat/chat_completion`;
  window.console.log('chat', url);

  const ctrl = new AbortController();

  const data = {
    model: 'qwen2.5-coder:7b-instruct',
    messages: [
      { role: 'user', content: '你好' },
      { role: 'assistant', content: '你好，我是人工智能大模型' },
      { role: 'user', content: userMsg },
    ],
    stream: true,
    temperature: 0.7,
    extra_body: {
      top_k: 3,
      score_threshold: 1.8,
      return_direct: false,
    },
    kb_name: kb.value,
  };

  window.console.log('======', data);

  isChatting.value = true;
  ctrl.abort();
  class FetchError extends Error {}

  await fetchEventSource(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
    mode: 'cors',
    signal: ctrl.signal,

    async onopen(response) {
      window.console.log('onopen', response);

      if (!response.ok) {
        throw new FetchError();
      }
    },

    onmessage(msg: any) {
      // return if no data
      if (!msg.data) return;

      window.console.log('onmessage', msg);

      let jsn = {} as any;
      try {
        jsn = JSON.parse(msg.data);
      } catch {
        window.console.log('parse chatchat msg failed', msg.data);
        throw new FetchError();
      }

      const doc_contents = [] as any[];
      let msg_content = '';

      // parse msg
      if (jsn.docs && jsn.docs.length > 0) {
        // docs
        const docMap = {} as any;

        jsn.docs.forEach((doc: any) => {
          const { pageId, pageTitle, pageType } = getDocLink(doc.trim());
          if (!docMap[pageId] && pageType === 'html') {
            // is link
            const doc_content = `[${pageTitle}](${wikiAddress}/pages/viewpage.action?pageId=${pageId})`;

            doc_contents.push(doc_content);
            docMap[pageId] = true;
          }
        });
      } else if (jsn.choices && jsn.choices.length > 0) {
        // msg
        jsn.choices?.forEach((choice: any) => {
          if (choice.delta?.content && choice.delta?.content !== '__BREAK__') {
            msg_content += choice.delta?.content;
          }
        });
      }

      // generate msg
      let docs = '';
      let content = '';
      if (doc_contents.length > 0) {
        docs = `  \n参考资料：\n1. ${doc_contents.join('  \n1. ')}`;
      } else if (msg_content.length > 0) {
        content = `${msg_content}`;
      }

      // create/update robot msg
      if (continueOnCurrMsg.value) {
        // append
        const index = getLatestRobotMsg(messages.value);
        if (index >= 0) {
          if (docs.length > 0)
            messages.value[index].docs = replaceLinkWithoutTitle(
              messages.value[index].docs + docs,
            );

          if (content.length > 0)
            messages.value[index].content = replaceLinkWithoutTitle(
              messages.value[index].content + content,
            );
        }
      } else {
        // create
        const currRobotMsg = {
          type: 'robot',
          name: humanName,
          avatar: humanAvatar,
          docs: docs.length > 0 ? replaceLinkWithoutTitle(docs) : '',
          content: content.length > 0 ? content : '',
        };
        // window.console.log('!!!!!!', currRobotMsg)
        messages.value.push(currRobotMsg);

        continueOnCurrMsg.value = true;
      }

      scroll();
    },

    onclose() {
      window.console.log(
        'onclose',
        messages.value.length > 0
          ? messages.value[messages.value.length - 1].content
          : 'empty',
      );
      isChatting.value = false;
      continueOnCurrMsg.value = false;
      ctrl.abort();
    },
    onerror(err: any) {
      window.console.log('onerror', err);
      isChatting.value = false;
      continueOnCurrMsg.value = false;

      ctrl.abort();
      throw err; // rethrow to stop retries
    },
  });
};

const keyDown = (event: any) => {
  window.console.log(event);

  if (historyIndex.value === -1 && histories.value.length === 0) {
    // fist time
    return;
  }

  if (event.keyCode === consts.keyCodeUp) {
    window.console.log('up');

    if (historyIndex.value === -1) {
      // fist time
      historyIndex.value = histories.value.length - 1;
      msg.value = histories.value[historyIndex.value];

      setSelectionRange(document.querySelector('#msgInput'), msg.value.length);

      return;
    }

    if (historyIndex.value > 0) {
      historyIndex.value--;
    }
    msg.value = histories.value[historyIndex.value];
  } else if (event.keyCode === consts.keyCodeDown) {
    window.console.log('keyDown', event);

    if (
      historyIndex.value === -1 || // fist time
      historyIndex.value === histories.value.length - 1
    ) {
      // is max
      historyIndex.value = -1;
      msg.value = '';
      return;
    }

    historyIndex.value++;
    msg.value = histories.value[historyIndex.value];
  }

  if (
    event.keyCode === consts.keyCodeUp ||
    event.keyCode === consts.keyCodeDown
  ) {
    setSelectionRange(document.querySelector('#msgInput'), msg.value.length);
  }
};

const initChatData = async () => {
  // const serverUrl = addSepIfNeeded(props.serverUrl);
};

const initHistory = async () => {
  histories.value = await getCache(CHAT_HISTORIES);
  if (!histories.value) histories.value = [];
  // if (histories.value.length > 0)
  //   msg.value = histories.value[histories.value.length - 1]
};

const showChat = ref(true);
const showOrNot = () => {
  showChat.value = !showChat.value;
};

const recall = (index: number) => {
  window.console.log('recall', index);
  if (index > messages.value.length - 1) {
    return;
  }

  const item = messages.value[index - 1];
  msg.value = item.content;
  send();
};

const copy = () => {
  window.console.log('copy');
  if (messages.value.length === 0 || !navigator.clipboard) {
    return;
  }

  navigator.clipboard.writeText(
    messages.value[messages.value.length - 1].content,
  );
  message.success('成功复制回复结果到剪贴板。');
};

const handleLinkClick = (event: any) => {
  window.console.log('handleLinkClick');

  const target = event.target;

  if (target.tagName.toLowerCase() === 'a' && target.getAttribute('href')) {
    if (!isUnderRobotMsg(target)) return true;

    event.preventDefault();

    const href = target.getAttribute('href');
    window.open(href, '_blank');
  }
};

onMounted(async () => {
  initHistory();
  initChatData();
  document.addEventListener('click', handleLinkClick);

  scroll();
});
onBeforeUnmount(async () => {
  document.removeEventListener('click', handleLinkClick);
});
</script>

<template>
  <div class="chatbot-main">
    <div
      v-if="!showChat"
      class="fix-action-open dp-link clear-both"
      title="开始聊天"
      @click="showOrNot"
    >
      <span class="open"></span>
    </div>

    <div v-if="showChat" class="chatbot-container">
      <div class="header">
        <div class="logo">
          <img src="/static/icon/chat-robot.png" />
        </div>

        <div class="label">ChatOPS</div>
      </div>

      <div id="chat-messages" class="messages">
        <template v-for="(item, index) in messages" :key="index">
          <div v-if="item.type === 'human'" class="chat-sender human">
            <div class="avatar-container">
              <div class="avatar"></div>
            </div>

            <div class="content">
              <span>{{ item.content }}</span>

              <span>{{ item.doc }}</span>

              <span
                v-if="isChatting && index === messages.length - 1"
                class="loading"
              >
                <img src="/static/icon/chat-loading.gif" />
              </span>
            </div>
          </div>

          <div v-if="item.type === 'robot'" class="chat-sender robot">
            <div class="avatar-container">
              <div class="avatar"></div>
            </div>

            <div class="content markdown-container">
              <span v-if="item.content">{{ item.content }}</span>
              <Markdown
                v-else
                :html="false"
                :linkify="true"
                :source="`${item.docs}\n\n${item.content}`"
              />
            </div>
            <div class="toolbar">
              <div class="call">
                <span class="dp-link-primary" @click="recall(index)">
                  重新生成
                </span>
              </div>

              <div class="copy dp-link" @click="copy">
                <img alt="copy" src="/static/logo.png" />
                复制
              </div>
            </div>
          </div>
        </template>
      </div>

      <div class="sender">
        <input
          id="msgInput"
          v-model="msg"
          autocomplete="off"
          class="input"
          @keydown="keyDown"
          @keyup.enter="send"
        />

        <span v-if="!isChatting" class="button dp-link" @click="send"></span>
        <span v-if="isChatting" class="button"></span>
      </div>

      <div class="actions">
        <slot name="actions"></slot>
      </div>
    </div>
  </div>
</template>

<style lang="less" src="./style.less" />
<style lang="less" src="./style-scoped.less" scoped />
