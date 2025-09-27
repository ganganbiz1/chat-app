"use client";

import { useState } from "react";
import { ChatRoomList } from "@/components/chat/chat-room-list";
import { ChatRoom } from "@/components/chat/chat-room";
import { CreateRoomModal } from "@/components/chat/create-room-modal";
import { Room } from "@/lib/api";

export default function ChatPage() {
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);

  const handleRoomSelect = (room: Room) => {
    setSelectedRoom(room);
  };

  const handleBackToList = () => {
    setSelectedRoom(null);
  };

  const handleCreateRoom = () => {
    setShowCreateModal(true);
  };

  const handleRoomCreated = (newRoom: Room) => {
    setShowCreateModal(false);
    setSelectedRoom(newRoom);
  };

  return (
    <div className="container mx-auto py-8 h-[calc(100vh-12rem)]">
      {selectedRoom ? (
        <ChatRoom room={selectedRoom} onBack={handleBackToList} />
      ) : (
        <ChatRoomList
          onRoomSelect={handleRoomSelect}
          onCreateRoom={handleCreateRoom}
        />
      )}

      <CreateRoomModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onRoomCreated={handleRoomCreated}
      />
    </div>
  );
}
