import { requestClient } from '#/api/request';

export async function clearKnowledgeBase(kb: string): Promise<any> {
  return requestClient.post(`/knowledgeBase/clearAll`, {}, { params: { kb } });
}
