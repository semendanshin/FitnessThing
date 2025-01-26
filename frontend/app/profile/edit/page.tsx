"use client";
import { Form } from "@nextui-org/form";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@nextui-org/input";
import { Avatar } from "@nextui-org/avatar";
import { DatePicker } from "@nextui-org/date-picker";
import { Button } from "@nextui-org/button";
import { CalendarDate } from "@internationalized/date";
import { toast } from "react-toastify";

import { PageHeader } from "@/components/page-header";
import { WorkoutUser } from "@/api/api.generated";
import { authApi, errUnauthorized } from "@/api/api";
import { Loading } from "@/components/loading";

function DataForm({ user }: { user: WorkoutUser }) {
  function protoToCalendarDate(s: string) {
    if (!s) {
      return null;
    }

    const date = new Date(s);

    return new CalendarDate(
      date.getUTCFullYear(),
      date.getUTCMonth() + 1,
      date.getUTCDate(),
    );
  }

  const [userDateOfBirth, setUserDateOfBirth] = useState(
    protoToCalendarDate(user.dateOfBirth!),
  );
  const [userFirstName, setUserFirstName] = useState(user.firstName);
  const [userLastName, setUserLastName] = useState(user.lastName);
  const [userWeight, setUserWeight] = useState(user.weight!);
  const [userHeight, setUserHeight] = useState(user.height!);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    await authApi.v1
      .userServiceUpdateUser({
        firstName: userFirstName ? userFirstName : undefined,
        lastName: userLastName ? userLastName : undefined,
        dateOfBirth: userDateOfBirth
          ? userDateOfBirth.toDate("UTC").toISOString()
          : undefined,
        weight: userWeight > 0 ? userWeight : undefined,
        height: userHeight > 0 ? userHeight : undefined,
      })
      .then(() => {
        toast.success("Данные успешно обновлены");
      })
      .catch((error) => {
        console.log(error);
        toast.error("Ошибка при обновлении данных");
      });
  }

  return (
    <Form className="flex flex-col pt-4 px-4 gap-4" onSubmit={handleSubmit}>
      <Input
        label="Имя"
        placeholder="Имя"
        value={userFirstName}
        onChange={(e) => setUserFirstName(e.target.value)}
      />
      <Input
        label="Фамилия"
        placeholder="Фамилия"
        value={userLastName}
        onChange={(e) => setUserLastName(e.target.value)}
      />
      <DatePicker
        label="Дата рождения"
        value={userDateOfBirth}
        onChange={setUserDateOfBirth}
      />
      <div className="flex flex-row gap-4 w-full">
        <Input
          className="w-1/2"
          label="Вес"
          type="number"
          value={String(userWeight)}
          onValueChange={(value) => setUserWeight(Number(value))}
        />
        <Input
          className="w-1/2"
          label="Рост"
          type="number"
          value={String(userHeight)}
          onValueChange={(value) => setUserHeight(Number(value))}
        />
      </div>
      <Button className="self-end w-full" color="success" type="submit">
        Сохранить
      </Button>
    </Form>
  );
}

function AvatarSection() {
  return (
    <div className="flex flex-col gap-4 items-center justify-around pt-4 px-4">
      <Avatar size="lg" src="https://i.pravatar.cc/300" />
    </div>
  );
}

export default function EditProfilePage() {
  const [user, setUser] = useState<WorkoutUser>();

  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const router = useRouter();

  async function fetchMe() {
    await authApi.v1
      .userServiceGetMe()
      .then((response) => {
        console.log(response.data);
        setUser(response.data.user!);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        throw error;
      });
  }

  async function fetchData() {
    setIsLoading(true);
    try {
      await fetchMe();
      setIsError(false);
    } catch {
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
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

  return (
    <div className="py-4 flex flex-col h-full">
      <PageHeader enableBackButton title="Редактировать профиль" />
      <AvatarSection />
      <DataForm user={user!} />
    </div>
  );
}
