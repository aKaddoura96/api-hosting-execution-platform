const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class APIClient {
  private getAuthHeader() {
    const token = localStorage.getItem('token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const authHeader = this.getAuthHeader();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(authHeader.Authorization ? { Authorization: authHeader.Authorization } : {}),
      ...(options.headers as Record<string, string> || {}),
    };

    const response = await fetch(url, { ...options, headers });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error || `HTTP ${response.status}`);
    }

    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  // Auth
  async signup(email: string, password: string, name: string, role: string) {
    return this.request('/api/v1/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ email, password, name, role }),
    });
  }

  async login(email: string, password: string) {
    return this.request('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async getMe() {
    return this.request('/api/v1/auth/me');
  }

  // APIs
  async getMyAPIs() {
    return this.request('/api/v1/apis');
  }

  async getPublicAPIs() {
    return this.request('/api/v1/marketplace/apis');
  }

  async getAPI(id: string) {
    return this.request(`/api/v1/apis/${id}`);
  }

  async createAPI(data: {
    name: string;
    description: string;
    version: string;
    runtime: string;
    visibility: string;
  }) {
    return this.request('/api/v1/apis', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async uploadCode(apiId: string, file: File) {
    const formData = new FormData();
    formData.append('code', file);

    const token = localStorage.getItem('token');
    const headers: Record<string, string> = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}/api/v1/apis/${apiId}/upload`, {
      method: 'POST',
      headers,
      body: formData,
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }

    return response.json();
  }

  async deleteAPI(id: string) {
    return this.request(`/api/v1/apis/${id}`, {
      method: 'DELETE',
    });
  }
}

export const api = new APIClient();
