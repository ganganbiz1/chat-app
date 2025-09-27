"use client";

import { useState, useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";
import { UserModal } from "./user-modal";
import { Room, Message, apiClient, getUserIdFromStorage, getUsernameFromStorage } from "@/lib/api";
import { ArrowLeft, Users, Settings } from "lucide-react";

interface ChatRoomProps {
  room: Room;
  onBack: () => void;
}

export function ChatRoom({ room, onBack }: ChatRoomProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showUserModal, setShowUserModal] = useState(false);
  const [userId] = useState(() => getUserIdFromStorage());
  const [username, setUsername] = useState(() => getUsernameFromStorage());
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  const loadMessages = async () => {
    try {
      setError(null);
      const recentMessages = await apiClient.getRecentMessages(room.id, 50);
      setMessages(recentMessages.reverse()); // 古いメッセージが上に来るように反転
    } catch (err) {
      setError(err instanceof Error ? err.message : "メッセージの読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
  };

  const sendMessage = async (content: string) => {
    try {
      const newMessage = await apiClient.createMessage(room.id, {
        user_id: userId,
        username,
        content,
      });

      // 新しいメッセージを配列の最後に追加
      setMessages(prev => [...prev, newMessage]);
    } catch (err) {
      console.error("メッセージの送信に失敗しました:", err);
      throw err;
    }
  };

  const pollMessages = async () => {
    try {
      const recentMessages = await apiClient.getRecentMessages(room.id, 50);
      const reversedMessages = recentMessages.reverse();
      
      // メッセージが変更されている場合のみ更新
      if (JSON.stringify(reversedMessages) !== JSON.stringify(messages)) {
        setMessages(reversedMessages);
      }
    } catch (err) {
      console.error("メッセージの更新に失敗しました:", err);
    }
  };

  useEffect(() => {
    loadMessages();

    // 5秒ごとにメッセージを更新（簡易的なリアルタイム機能）
    intervalRef.current = setInterval(pollMessages, 5000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [room.id]);

  // usernameが変更されたときにローカルストレージを更新
  useEffect(() => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('chat_username', username);
    }
  }, [username]);

  const handleUsernameChange = (newUsername: string) => {
    setUsername(newUsername);
    setShowUserModal(false);
  };

  return (
    <div className="h-full flex flex-col">
      {/* ヘッダー */}
      <Card className="mb-4">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={onBack}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <div>
                <CardTitle className="text-xl">{room.name}</CardTitle>
                {room.description && (
                  <p className="text-sm text-gray-600 mt-1">{room.description}</p>
                )}
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowUserModal(true)}
              >
                <Users className="h-4 w-4 mr-2" />
                {username}
              </Button>
            </div>
          </div>
        </CardHeader>
      </Card>

      {/* メッセージエリア */}
      <div className="flex-1 flex flex-col min-h-0">
        <Card className="flex-1 flex flex-col min-h-0">
          <CardContent className="flex-1 flex flex-col p-0 min-h-0">
            <MessageList
              messages={messages}
              loading={loading}
              error={error}
              currentUserId={userId}
              onRetry={loadMessages}
            />
            <MessageInput
              onSendMessage={sendMessage}
              disabled={loading}
            />
          </CardContent>
        </Card>
      </div>

      {/* ユーザー設定モーダル */}
      <UserModal
        isOpen={showUserModal}
        onClose={() => setShowUserModal(false)}
        currentUsername={username}
        onUsernameChange={handleUsernameChange}
      />
    </div>
  );
}

