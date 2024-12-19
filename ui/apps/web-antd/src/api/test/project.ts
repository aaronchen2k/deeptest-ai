import { requestClient } from '#/api/request';

export async function listProjectApi(data: any) {
  return requestClient.post('/projects/query', data);
}
