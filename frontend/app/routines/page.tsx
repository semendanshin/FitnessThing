"use client";
import { Button } from "@nextui-org/button";
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import {
  Form,
  Input,
  Link,
  Modal,
  ModalContent,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { PlusIcon } from "@/config/icons";
import { PageHeader } from "@/components/page-header";
import { Loading } from "@/components/loading";
import { authApi, errUnauthorized } from "@/api/api";

export default function RoutinesPage() {
  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const { isOpen, onOpen, onOpenChange } = useDisclosure();

  const [routines, setRoutines] = useState<any[]>([]);

  const router = useRouter();

  async function fetchData() {
    setIsLoading(true);
    authApi.v1
      .routineServiceGetRoutines()
      .then((response) => {
        console.log(response.data);
        setIsError(false);
        setRoutines(response.data.routines!);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized) {
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

  return (
    <>
      <div className="py-4 flex-grow">
        <PageHeader enableBackButton title="Шаблоны" />
        <div className="grid grid-cols-1 gap-4 p-4">
          {routines.map((routine) => (
            <Card
              key={routine.id}
              fullWidth
              as={Link}
              className="flex flex-row flex-grow p-2 gap-4 justify-between"
              href={`/routines/${routine.id}`}
              shadow="sm"
            >
              <div className="flex flex-col items-start justify-between p-2">
                <CardHeader className="p-0">
                  <p className="text-m font-bold">{routine.name}</p>
                </CardHeader>
                <CardBody className="p-0">
                  {routine.description ? (
                    <div className="py-2">
                      <p className="text-xs text-gray-400/80">
                        {routine.description}
                      </p>
                    </div>
                  ) : null}
                  <p className="text-xs text-gray-500">
                    {Math.ceil(Math.random() * (10 - 5) + 5)} упражнений
                  </p>
                </CardBody>
              </div>
            </Card>
          ))}
          <Button color="primary" onPress={onOpen}>
            <PlusIcon className="w-4 h-4" />
            Добавить рутину
          </Button>
        </div>
      </div>
      <Modal
        isOpen={isOpen}
        placement="center"
        size="xs"
        onClose={onOpenChange}
      >
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader>
                <h2>Добавить рутину</h2>
              </ModalHeader>
              <Form
                className="inline-block text-center justify-center w-full max-w-lg p-4"
                validationBehavior="native"
                onSubmit={async (e: React.FormEvent<HTMLFormElement>) => {
                  e.preventDefault();
                  const data = Object.fromEntries(
                    new FormData(e.currentTarget),
                  );

                  console.log(data);
                  await authApi.v1
                    .routineServiceCreateRoutine({
                      name: data.name.toString(),
                      description: data.description.toString(),
                    })
                    .then((response) => {
                      console.log(response.data);
                      router.push(`/routines/${response.data.routine?.id}`);
                    })
                    .catch((error) => {
                      console.log(error);

                      return error;
                    });
                }}
              >
                <div className="flex flex-col items-center justify-center gap-4 py-4">
                  <Input
                    autoFocus
                    isRequired
                    label="Название"
                    labelPlacement="outside"
                    name="name"
                    placeholder="Название"
                    type="text"
                  />
                  <Input
                    label="Описание"
                    labelPlacement="outside"
                    name="description"
                    placeholder="Описание"
                    type="text"
                  />
                  <Button color="primary" type="submit">
                    Добавить
                  </Button>
                </div>
              </Form>
            </>
          )}
        </ModalContent>
      </Modal>
    </>
  );
}
