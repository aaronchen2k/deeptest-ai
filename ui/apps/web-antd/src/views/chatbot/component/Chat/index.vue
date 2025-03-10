<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';

import { useAppConfig } from '@vben/hooks';
import { useUserStore } from '@vben/stores';
import { addSepIfNeeded } from '@vben/utils';

import { fetchEventSource } from '@microsoft/fetch-event-source';
import { message } from 'ant-design-vue';
import Markdown from 'vue3-markdown-it';

import consts from '#/config/constant';
import { getCache, setCache } from '#/utils/cache-local';

import Uploder from '../Uploder/index.vue';
import {
  getDocLink,
  getLatestRobotMsg,
  getMaterialIdInDocName,
  isUnderRobotMsg,
  replaceImageUrl,
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
  serverUrl: {
    type: String,
    default: 'http://127.0.0.1:9085/api/v1',
    required: true,
  },
});

const userStore = useUserStore();

const { docRepoUrl } = useAppConfig(import.meta.env, import.meta.env.PROD);
const wakeUpWord = '小深';
const humanName = 'Albert';
const humanAvatar = '/static/chat-einstein.png';

const CHAT_HISTORIES = 'chat_history_key';
const histories = ref([] as any[]);
const historyIndex = ref(-1);

const msg = ref('');
const conversationId = ref('' as string);
const isChatting = ref(false);
const continueOnCurrMsg = ref(false);

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
    content: '你好，有什么可以帮助你的吗？',
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
  conversationId.value = '';
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
  const url = `${serverUrl}chatbot/chat`;
  window.console.log('chat', url);

  const ctrl = new AbortController();

  const data = {
    user: `${userStore.userInfo?.id}`,
    query: userMsg,
    response_mode: 'streaming',
    temperature: 0.7,
    inputs: {},
    has_thoughts: true,
  };
  window.console.log('====== request data:', data);

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
      window.console.log('------ onopen', response);

      if (!response.ok) {
        throw new FetchError();
      }
    },

    onmessage(msg: any) {
      if (!msg.data) return;
      window.console.log('------ onmessage', msg);

      let jsn = {} as any;
      try {
        jsn = JSON.parse(msg.data);
      } catch {
        window.console.log('parse chat response failed', msg.data);
        throw new FetchError();
      }

      const doc_contents = [] as any[];
      let msg_content = '';

      // drop if not same conversation
      if (
        conversationId.value !== '' &&
        jsn.conversation_id !== '' &&
        conversationId.value !== jsn.conversation_id
      ) {
        return;
      }

      if (conversationId.value === '' && jsn.conversation_id !== '') {
        conversationId.value = jsn.conversation_id;
      }

      if (jsn.answer) {
        // answer
        msg_content += jsn.answer;
      } else if (jsn.metadata?.retriever_resources) {
        // data
        const docMap = {} as any;

        jsn.metadata?.retriever_resources.forEach((res: any) => {
          const { docId, docName, docType } = getDocLink(res);
          if (!docMap[docId] && docType === 'upload_file') {
            const { materialId, newDocName } = getMaterialIdInDocName(docName);
            const url = `${docRepoUrl}/${materialId}/${docName}`; // important
            const doc_content = `[${newDocName}](${url})`;

            doc_contents.push(doc_content);
            docMap[docId] = true;
          }
        });
      }

      // generate msg
      let docs = '';
      let content = '';
      if (doc_contents.length > 0) {
        docs = `  \n *参考资料：* \n1. ${doc_contents.join('  \n1. ')}\n`;
      } else if (msg_content.length > 0) {
        content = `${msg_content}`;
      }

      // create/update robot msg
      if (continueOnCurrMsg.value) {
        // append
        const index = getLatestRobotMsg(messages.value);
        if (index >= 0) {
          if (content.length > 0)
            // add http base url for image
            // ![剧照](7/./img/stills_02.jpg "剧照")
            messages.value[index].content = replaceImageUrl(
              messages.value[index].content + content,
              docRepoUrl,
            );

          if (docs.length > 0)
            // add http base url for doc
            // [百年孤独_8.md](http://localhost:9085/upload/public/knowledge_base/百年孤独_8.md)
            messages.value[index].docs = messages.value[index].docs + docs;
          // messages.value[index].docs = replaceDocUrl(
          //   messages.value[index].docs + docs,
          //   docRepoUrl,
          // );
          window.console.log('====== docs like', messages.value[index].docs);

          window.console.log('!!!', messages.value[index]);
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

function getRobotMsg(item: any) {
  return item.content || item.docs ? `${item.content}\n\n${item.docs}` : '...';
}

onMounted(async () => {
  initHistory();
  document.addEventListener('click', handleLinkClick);

  scroll();
});
onBeforeUnmount(async () => {
  document.removeEventListener('click', handleLinkClick);
});
</script>

<template>
  <div class="chat-main">
    <div
      v-if="!showChat"
      class="open-window-btn dp-link clear-both"
      title="开始聊天"
      @click="showOrNot"
    >
      <span class="open"></span>
    </div>

    <div v-if="showChat" class="chatbot-container">
      <div class="header">
        <div class="logo">
          <img src="/static/icon/chat-logo.png" />
        </div>

        <div class="label">ChatOPS</div>

        <div class="contrl">
          <div class="select-kb-wrapper"></div>
          <div class="action dp-link" @click="showOrNot">
            <span class="close"></span>
          </div>
        </div>
      </div>

      <div id="chat-messages" class="messages">
        <template v-for="(item, index) in messages" :key="index">
          <div v-if="item.type === 'human'" class="chat-record human">
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

          <div v-if="item.type === 'robot'" class="chat-record robot">
            <div class="avatar-container">
              <div class="avatar"></div>
            </div>

            <div class="content markdown-container">
              <Markdown
                :html="false"
                :linkify="true"
                :source="getRobotMsg(item)"
              />
            </div>
            <div style="clear: both"></div>
            <div class="toolbar">
              <div class="call">
                <span class="dp-link-primary" @click="recall(index)">
                  重新生成
                </span>
              </div>

              <div class="copy dp-link" @click="copy">复制</div>
            </div>
          </div>
        </template>
      </div>

      <div class="sender">
        <textarea
          id="msgInput"
          v-model="msg"
          autocomplete="off"
          class="input"
          placeholder="可使用上下键切换历史聊天记录"
          rows="2"
          @keydown="keyDown"
          @keyup.enter="send"
        ></textarea>

        <span v-if="!isChatting" class="button dp-link" @click="send"></span>
        <span v-if="isChatting" class="button"></span>
      </div>

      <div class="actions">
        <div class="uploader-container">
          <Uploder />
        </div>
        <slot name="actions"></slot>
      </div>
    </div>
  </div>
</template>

<style lang="less" src="./style.less" />
<style lang="less" src="./style-scoped.less" scoped />
