import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function HomePage() {
  return (
    <div className="container mx-auto py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Chat App へようこそ
        </h1>
        <p className="text-lg text-gray-600">
          リアルタイムでチャットを楽しもう
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              📢 お知らせ
            </CardTitle>
            <CardDescription>
              最新のお知らせとアップデート情報
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">
              お知らせページでアプリの最新情報をチェックしましょう。
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              💬 チャット
            </CardTitle>
            <CardDescription>
              リアルタイムチャットルーム
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">
              チャットページで他のユーザーとリアルタイムでコミュニケーションを取りましょう。
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}