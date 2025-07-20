import { GetTokenSilentlyOptions } from '@auth0/auth0-react';

type GetAccessToken = (options?: GetTokenSilentlyOptions) => Promise<string>;

export const createApiClient = (getAccessToken: GetAccessToken) => {
  const request = async (
    url: string,
    method: 'GET' | 'POST' | 'PATCH' | 'DELETE',
    body?: any
  ) => {
    try {
      const token = await getAccessToken({
        authorizationParams: {

        },
      });
      
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      };

      const options: RequestInit = {
        method, 
        headers, 
        body: body ? JSON.stringify(body) : undefined,
      };

      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}${url}`, options);

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      return response.json();
    } catch (error) {
      console.error(`API request failed: ${method} ${url}`, error);
      throw error;
    }
  };

  return {
    get: (url: string) => request(url, 'GET'),
    post: (url: string, body: any) => request(url, 'POST', body),
    patch: (url: string, body: any) => request(url, 'PATCH', body),
 };
};
