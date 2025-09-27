// API client for chat app backend

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Types
export interface Room {
  id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: string;
  room_id: string;
  user_id: string;
  username: string;
  content: string;
  created_at: string;
}

export interface CreateRoomRequest {
  name: string;
  description: string;
}

export interface CreateMessageRequest {
  user_id: string;
  username: string;
  content: string;
}

export interface PaginationResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface ErrorResponse {
  error: string;
  message?: string;
  code?: number;
}

// API functions
class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData: ErrorResponse = await response.json().catch(() => ({
          error: 'Network error',
          message: `HTTP ${response.status}: ${response.statusText}`,
        }));
        throw new Error(errorData.message || errorData.error);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Room API methods
  async createRoom(data: CreateRoomRequest): Promise<Room> {
    const response = await this.request<ApiResponse<Room>>('/api/rooms', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data;
  }

  async getRooms(page = 1, limit = 20): Promise<PaginationResponse<Room>> {
    return await this.request<PaginationResponse<Room>>(
      `/api/rooms?page=${page}&limit=${limit}`
    );
  }

  async getRoom(id: string): Promise<Room> {
    const response = await this.request<ApiResponse<Room>>(`/api/rooms/${id}`);
    return response.data;
  }

  async updateRoom(id: string, data: Partial<CreateRoomRequest>): Promise<Room> {
    const response = await this.request<ApiResponse<Room>>(`/api/rooms/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return response.data;
  }

  async deleteRoom(id: string): Promise<void> {
    await this.request(`/api/rooms/${id}`, {
      method: 'DELETE',
    });
  }

  // Message API methods
  async createMessage(roomId: string, data: CreateMessageRequest): Promise<Message> {
    const response = await this.request<ApiResponse<Message>>(
      `/api/rooms/${roomId}/messages`,
      {
        method: 'POST',
        body: JSON.stringify(data),
      }
    );
    return response.data;
  }

  async getRoomMessages(
    roomId: string,
    page = 1,
    limit = 50
  ): Promise<PaginationResponse<Message>> {
    return await this.request<PaginationResponse<Message>>(
      `/api/rooms/${roomId}/messages?page=${page}&limit=${limit}`
    );
  }

  async getRecentMessages(roomId: string, limit = 20): Promise<Message[]> {
    const response = await this.request<ApiResponse<Message[]>>(
      `/api/rooms/${roomId}/messages/recent?limit=${limit}`
    );
    return response.data;
  }

  async getMessage(id: string): Promise<Message> {
    const response = await this.request<ApiResponse<Message>>(`/api/messages/${id}`);
    return response.data;
  }

  async deleteMessage(id: string): Promise<void> {
    await this.request(`/api/messages/${id}`, {
      method: 'DELETE',
    });
  }

  async getUserMessages(
    userId: string,
    page = 1,
    limit = 50
  ): Promise<PaginationResponse<Message>> {
    return await this.request<PaginationResponse<Message>>(
      `/api/users/${userId}/messages?page=${page}&limit=${limit}`
    );
  }

  // Health check
  async healthCheck(): Promise<{ status: string; service: string }> {
    return await this.request('/health');
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Utility functions
export const formatMessageTime = (timestamp: string): string => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60));

  if (diffInMinutes < 1) {
    return 'たった今';
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分前`;
  } else if (diffInMinutes < 1440) {
    const hours = Math.floor(diffInMinutes / 60);
    return `${hours}時間前`;
  } else {
    return date.toLocaleDateString('ja-JP', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
};

export const generateUserId = (): string => {
  return `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
};

export const getUsernameFromStorage = (): string => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('chat_username') || 'ゲスト';
  }
  return 'ゲスト';
};

export const setUsernameToStorage = (username: string): void => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('chat_username', username);
  }
};

export const getUserIdFromStorage = (): string => {
  if (typeof window !== 'undefined') {
    let userId = localStorage.getItem('chat_user_id');
    if (!userId) {
      userId = generateUserId();
      localStorage.setItem('chat_user_id', userId);
    }
    return userId;
  }
  return generateUserId();
};

