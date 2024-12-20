import { requestClient } from '#/api/request';

export async function listProjectApi(data: any) {
  return requestClient.post('/projects/query', data);
}

export async function getProjectApi(id: number) {
  return requestClient.get(`/projects/${id}`);
}

export async function saveProjectApi(data: any) {
  const method = data.id ? 'put' : 'post';
  return requestClient[method](`/projects`, data);
}
