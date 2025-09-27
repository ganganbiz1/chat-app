import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

// サンプルお知らせデータ
const announcements = [
  {
    id: 1,
    title: "チャットアプリがリリースされました！",
    content: "新しいリアルタイムチャットアプリが正式にリリースされました。ぜひお試しください。",
    date: "2024-09-27",
  },
  {
    id: 2,
    title: "メンテナンスのお知らせ",
    content: "システムメンテナンスのため、9月28日 AM 2:00-4:00の間、一時的にサービスを停止いたします。",
    date: "2024-09-26",
  },
  {
    id: 3,
    title: "新機能：ファイル共有機能を追加予定",
    content: "次回のアップデートで、チャットでファイルを共有できる機能を追加予定です。",
    date: "2024-09-25",
  },
];

export default function AnnouncementsPage() {
  return (
    <div className="container mx-auto py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2 flex items-center gap-3">
          📢 お知らせ
        </h1>
        <p className="text-lg text-gray-600">
          アプリの最新情報とお知らせ
        </p>
      </div>

      <div className="space-y-6">
        {announcements.map((announcement) => (
          <Card key={announcement.id}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>{announcement.title}</CardTitle>
                <span className="text-sm text-gray-500">
                  {announcement.date}
                </span>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-gray-700 leading-relaxed">
                {announcement.content}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
