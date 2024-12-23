import { requestClient } from '#/api/request';

export async function loadCaseApi() {
  return requestClient.post('/cases/load', {});
}

export async function getCaseApi(id: number) {
  return requestClient.get(`/cases/${id}`);
}

export async function saveCaseApi(data: any) {
  const method = data.id ? 'put' : 'post';
  return requestClient[method](`/cases`, data);
}
