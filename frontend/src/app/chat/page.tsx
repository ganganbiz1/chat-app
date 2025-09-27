import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

// サンプルチャットルームデータ
const chatRooms = [
  {
    id: 1,
    name: "一般チャット",
    description: "誰でも参加できる一般的なチャットルームです",
    memberCount: 12,
    lastMessage: "お疲れさまです！",
    lastMessageTime: "10分前",
  },
  {
    id: 2,
    name: "技術討論",
    description: "プログラミングや技術に関する話題を話しましょう",
    memberCount: 8,
    lastMessage: "Next.jsの新機能について",
    lastMessageTime: "30分前",
  },
  {
    id: 3,
    name: "雑談ルーム",
    description: "自由に雑談を楽しむルームです",
    memberCount: 15,
    lastMessage: "今日はいい天気ですね",
    lastMessageTime: "1時間前",
  },
];

export default function ChatPage() {
  return (
    <div className="container mx-auto py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2 flex items-center gap-3">
          💬 チャット
        </h1>
        <p className="text-lg text-gray-600">
          チャットルームを選択して会話を始めましょう
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {chatRooms.map((room) => (
          <Card key={room.id} className="hover:shadow-md transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                {room.name}
                <span className="text-sm font-normal text-gray-500">
                  {room.memberCount}人
                </span>
              </CardTitle>
              <CardDescription>
                {room.description}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="text-sm text-gray-600">
                  <p className="font-medium">最新メッセージ:</p>
                  <p className="italic">"{room.lastMessage}"</p>
                  <p className="text-xs text-gray-500 mt-1">
                    {room.lastMessageTime}
                  </p>
                </div>
                <Button className="w-full">
                  ルームに参加
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="mt-8 text-center">
        <Button variant="outline" size="lg">
          新しいルームを作成
        </Button>
      </div>
    </div>
  );
}
