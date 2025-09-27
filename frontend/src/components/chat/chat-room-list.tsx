"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Room, apiClient, formatMessageTime } from "@/lib/api";
import { Loader2, Plus, MessageCircle } from "lucide-react";

interface ChatRoomListProps {
  onRoomSelect: (room: Room) => void;
  onCreateRoom: () => void;
}

export function ChatRoomList({ onRoomSelect, onCreateRoom }: ChatRoomListProps) {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadRooms = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.getRooms(1, 20);
      setRooms(response.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ルームの読み込みに失敗しました');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadRooms();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="flex items-center gap-2">
          <Loader2 className="h-6 w-6 animate-spin" />
          <span>ルームを読み込み中...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <div className="text-red-600 mb-4">{error}</div>
        <Button onClick={loadRooms} variant="outline">
          再試行
        </Button>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2 flex items-center gap-3">
          <MessageCircle className="h-8 w-8" />
          チャット
        </h1>
        <p className="text-lg text-gray-600">
          チャットルームを選択して会話を始めましょう
        </p>
      </div>

      {rooms.length === 0 ? (
        <div className="text-center py-12">
          <MessageCircle className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            まだルームがありません
          </h3>
          <p className="text-gray-600 mb-6">
            最初のチャットルームを作成して会話を始めましょう
          </p>
          <Button onClick={onCreateRoom} size="lg">
            <Plus className="h-5 w-5 mr-2" />
            新しいルームを作成
          </Button>
        </div>
      ) : (
        <>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {rooms.map((room) => (
              <Card key={room.id} className="hover:shadow-md transition-shadow cursor-pointer">
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    {room.name}
                    <span className="text-sm font-normal text-gray-500">
                      {/* TODO: Add online user count */}
                    </span>
                  </CardTitle>
                  <CardDescription>
                    {room.description || "説明なし"}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="text-sm text-gray-600">
                      <p className="font-medium">作成日時:</p>
                      <p className="text-xs text-gray-500">
                        {formatMessageTime(room.created_at)}
                      </p>
                    </div>
                    <Button 
                      className="w-full" 
                      onClick={() => onRoomSelect(room)}
                    >
                      ルームに参加
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          <div className="mt-8 text-center">
            <Button onClick={onCreateRoom} variant="outline" size="lg">
              <Plus className="h-5 w-5 mr-2" />
              新しいルームを作成
            </Button>
          </div>
        </>
      )}
    </div>
  );
}

