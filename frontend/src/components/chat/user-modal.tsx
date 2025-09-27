"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { X, User } from "lucide-react";

interface UserModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentUsername: string;
  onUsernameChange: (username: string) => void;
}

export function UserModal({
  isOpen,
  onClose,
  currentUsername,
  onUsernameChange,
}: UserModalProps) {
  const [username, setUsername] = useState(currentUsername);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const trimmedUsername = username.trim();
    if (!trimmedUsername) {
      setError("ユーザー名を入力してください");
      return;
    }

    if (trimmedUsername.length > 255) {
      setError("ユーザー名は255文字以内で入力してください");
      return;
    }

    onUsernameChange(trimmedUsername);
    setError(null);
  };

  const handleClose = () => {
    setUsername(currentUsername);
    setError(null);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <User className="h-5 w-5" />
                <div>
                  <CardTitle>ユーザー設定</CardTitle>
                  <CardDescription>
                    チャットで表示される名前を設定してください
                  </CardDescription>
                </div>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleClose}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
                  ユーザー名 *
                </label>
                <input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="例: 太郎"
                  maxLength={255}
                  required
                />
                <p className="text-xs text-gray-500 mt-1">
                  この名前がチャットメッセージで表示されます
                </p>
              </div>

              {error && (
                <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md">
                  {error}
                </div>
              )}

              <div className="flex gap-3 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleClose}
                  className="flex-1"
                >
                  キャンセル
                </Button>
                <Button
                  type="submit"
                  disabled={!username.trim()}
                  className="flex-1"
                >
                  保存
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

