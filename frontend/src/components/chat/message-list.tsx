"use client";

import { useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Message, formatMessageTime } from "@/lib/api";
import { MessageCircle, Loader2, AlertCircle } from "lucide-react";

interface MessageListProps {
  messages: Message[];
  loading: boolean;
  error: string | null;
  currentUserId: string;
  onRetry: () => void;
}

interface MessageBubbleProps {
  message: Message;
  isOwn: boolean;
}

function MessageBubble({ message, isOwn }: MessageBubbleProps) {
  return (
    <div className={`flex ${isOwn ? "justify-end" : "justify-start"} mb-4`}>
      <div className={`max-w-xs lg:max-w-md ${isOwn ? "order-2" : "order-1"}`}>
        <div
          className={`rounded-lg px-4 py-2 ${
            isOwn
              ? "bg-blue-500 text-white"
              : "bg-gray-100 text-gray-900"
          }`}
        >
          {!isOwn && (
            <div className="text-xs font-semibold mb-1 opacity-75">
              {message.username}
            </div>
          )}
          <div className="text-sm whitespace-pre-wrap break-words">
            {message.content}
          </div>
          <div
            className={`text-xs mt-1 ${
              isOwn ? "text-blue-100" : "text-gray-500"
            }`}
          >
            {formatMessageTime(message.created_at)}
          </div>
        </div>
      </div>
    </div>
  );
}

export function MessageList({
  messages,
  loading,
  error,
  currentUserId,
  onRetry,
}: MessageListProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="flex items-center gap-2 text-gray-600">
          <Loader2 className="h-6 w-6 animate-spin" />
          <span>メッセージを読み込み中...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="text-center">
          <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <div className="text-red-600 mb-4">{error}</div>
          <Button onClick={onRetry} variant="outline">
            再試行
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col min-h-0">
      <div
        ref={containerRef}
        className="flex-1 overflow-y-auto p-4 space-y-2"
        style={{ maxHeight: "calc(100vh - 300px)" }}
      >
        {messages.length === 0 ? (
          <div className="flex-1 flex items-center justify-center text-center py-12">
            <div>
              <MessageCircle className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                まだメッセージがありません
              </h3>
              <p className="text-gray-600">
                最初のメッセージを送信して会話を始めましょう
              </p>
            </div>
          </div>
        ) : (
          <>
            {messages.map((message) => (
              <MessageBubble
                key={message.id}
                message={message}
                isOwn={message.user_id === currentUserId}
              />
            ))}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>
    </div>
  );
}

