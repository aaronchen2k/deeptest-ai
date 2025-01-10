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

export async function remove(id: number): Promise<any> {
  return requestClient.delete(`cases/${id}`);
}

export async function createDirApi(data: any) {
  return requestClient.post(`/cases`, data);
}

export async function moveCaseApi(data: any) {
  return requestClient.post(`/cases/move`, data);
}
