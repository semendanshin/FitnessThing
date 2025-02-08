"use client";
import { useTheme } from "next-themes";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Avatar } from "@nextui-org/avatar";
import { Button } from "@nextui-org/button";
import { Card } from "@nextui-org/card";
import { Divider } from "@nextui-org/divider";

import { authApi, errUnauthorized } from "@/api/api";
import { Loading } from "@/components/loading";
import {
  ChevronRightIcon,
  EditIcon,
  ListIcon,
  TrophyIcon,
} from "@/config/icons";

export default function ProfilePage() {
  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const [user, setUser] = useState<any>({});

  const router = useRouter();
  const { theme, setTheme } = useTheme();

  async function fetchData() {
    setIsLoading(true);
    authApi.v1
      .userServiceGetMe()
      .then((response) => {
        console.log(response.data);
        setIsError(false);
        setUser(response.data.user!);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        setIsError(true);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  useEffect(() => {
    fetchData();
  }, []);

  if (isLoading) {
    return <Loading />;
  }

  if (isError) {
    return (
      <div className="p-4">
        <h2 className="text-lg text-red-500">Ошибка при загрузке данных</h2>
        <p>Проверьте соединение с сервером или обновите страницу.</p>
      </div>
    );
  }

  async function handleLogout() {
    await authApi.v1.authServiceLogout({
      refreshToken: localStorage.getItem("refreshToken")!,
    });
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    router.push("/auth/login");
  }

  function SubPageButton({
    href,
    label,
    icon,
  }: {
    href: string;
    label: string;
    icon?: React.ReactNode;
  }) {
    return (
      <Card isPressable className="p-3" onPress={() => router.push(href)}>
        <div className="flex flex-row justify-between items-center w-full">
          <div className="flex flex-row items-center gap-2">
            {icon}
            <p>{label}</p>
          </div>
          <ChevronRightIcon className="w-4 h-4" fill="currentColor" />
        </div>
      </Card>
    );
  }

  return (
    <div className="p-4 flex-grow gap-4">
      <div className="grid grid-cols-1 gap-4 py-4">
        <div className="flex flex-col gap-4 items-center justify-around">
          <Avatar size="lg" src={user.profilePictureUrl} />
          <h2 className="text-2xl font-bold">
            {user.firstName} {user.lastName}
          </h2>
        </div>
        <Divider />
        {/* Кнопочки история трениировок и еще какие-то которые я не придумал */}
        <SubPageButton
          href="/profile/edit"
          icon={<EditIcon className="w-4 h-4" fill="currentColor" />}
          label="Редактировать профиль"
        />
        <SubPageButton
          href="/profile/workouts"
          icon={<ListIcon className="w-4 h-4" fill="currentColor" />}
          label="История тренировок"
        />
        <SubPageButton
          href="/profile/records"
          icon={<TrophyIcon className="w-4 h-4" fill="currentColor" />}
          label="Рекорды"
        />
        {/*  Light and dark mode switch */}
        <Divider />
        <Button
          color="warning"
          onPress={() => {
            setTheme(theme === "dark" ? "light" : "dark");
          }}
        >
          Сменить тему
        </Button>
        <Divider />
        <Button color="danger" onPress={handleLogout}>
          Выйти
        </Button>
      </div>
    </div>
  );
}
