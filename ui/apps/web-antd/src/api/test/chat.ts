import { requestClient } from '#/api/request';

export async function listKnowledgeBase(): Promise<any> {
  return requestClient.get(`/list_knowledge_base`);
}
