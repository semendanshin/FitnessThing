"use client";

import { useRouter } from "next/navigation";
import { use, useEffect, useRef, useState } from "react";
import { toast } from "react-toastify";
import { Button } from "@nextui-org/button";
import { Card, CardBody } from "@nextui-org/card";
import { DropdownItem } from "@nextui-org/dropdown";
import { Input } from "@nextui-org/input";

import { Loading } from "@/components/loading";
import { authApi, errUnauthorized } from "@/api/api";
import {
  WorkoutExerciseInstanceDetails,
  WorkoutSet,
  WorkoutSetType,
} from "@/api/api.generated";
import { PageHeader } from "@/components/page-header";
import { PlusIcon, RepeatIcon, TrashCanIcon } from "@/config/icons";

export default function ExerciseInstancePage({
  params,
}: {
  params: Promise<{ id: string; exerciseInstanceId: string }>;
}) {
  const { exerciseInstanceId, id } = use(params);

  const [exerciseInstanceDetails, setExerciseInstanceDetails] =
    useState<WorkoutExerciseInstanceDetails>({});

  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const router = useRouter();

  async function fetchExerciseInstanceDetails() {
    await authApi.v1
      .routineServiceGetExerciseInstanceDetails(id, exerciseInstanceId)
      .then((response) => {
        console.log(response.data);
        setExerciseInstanceDetails(response.data.exerciseInstanceDetails!);
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
      await fetchExerciseInstanceDetails();
      setIsError(false);
    } catch (error) {
      console.log(error);
      toast.error("Failed to fetch workout details");
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

  function SetCard({ num, set }: { num: number; set: WorkoutSet }) {
    const [value, setValue] = useState(set.reps!);
    const timeoutRef = useRef<NodeJS.Timeout>();
    const valueRef = useRef(value); // Храним актуальное значение

    useEffect(() => {
      valueRef.current = value;
    }, [value]);

    useEffect(() => {
      return () => {
        if (timeoutRef.current) clearTimeout(timeoutRef.current);
      };
    }, []);

    async function updateSet() {
      try {
        await authApi.v1
          .routineServiceUpdateSetInExerciseInstance(
            id,
            exerciseInstanceId,
            set.id!,
            {
              reps: valueRef.current, // Используем актуальное значение из ref
              setType: WorkoutSetType.SET_TYPE_REPS,
            },
          )
          .then((response) => {
            console.log(response.data);
            fetchExerciseInstanceDetails();
          });
      } catch (error) {
        console.log(error);
        toast.error("Не удалось обновить подход");
      }
    }

    async function onDelete() {
      await authApi.v1
        .routineServiceRemoveSetFromExerciseInstance(
          id,
          exerciseInstanceId,
          set.id!,
        )
        .then((response) => {
          console.log(response.data);
          fetchExerciseInstanceDetails();
        })
        .catch((error) => {
          console.log(error);
          toast.error("Не удалось удалить подход");
        });
    }

    function onChange(e: React.ChangeEvent<HTMLInputElement>) {
      const newValue = parseInt(e.target.value) || 0;

      setValue(newValue); // Мгновенное обновление UI

      // Сбрасываем предыдущий таймер
      if (timeoutRef.current) clearTimeout(timeoutRef.current);

      // Устанавливаем новый таймер
      timeoutRef.current = setTimeout(() => {
        updateSet();
      }, 1000);
    }

    return (
      <Card>
        <CardBody className="flex flex-col gap-4">
          <div className="flex flex-row justify-between items-center">
            <div className="flex flex-row gap-2 items-center">
              <h2 className="text-sm font-semibold w-4">{num}.</h2>
              <div className="bg-default-100 rounded-small p-2">
                {(() => {
                  switch (set.setType) {
                    case WorkoutSetType.SET_TYPE_REPS:
                      return <RepeatIcon className="w-4 h-4" />;
                    default:
                      return <div className="w-4 h-4" />;
                  }
                })()}
              </div>
            </div>
            <div className="flex flex-row gap-6 items-center">
              <div className="flex flex-row gap-1 items-center">
                {(() => {
                  switch (set.setType) {
                    case WorkoutSetType.SET_TYPE_REPS:
                      return (
                        <div className="flex flex-row gap-1 items-center">
                          <Input
                            className="w-16"
                            size="sm"
                            type="number"
                            value={value.toString()}
                            onChange={onChange}
                          />
                          <p>раз</p>
                        </div>
                      );
                    default:
                      return null;
                  }
                })()}
              </div>
              <Button
                isIconOnly
                className="min-w-fit h-fit py-2 px-2"
                color="danger"
                size="sm"
                onPress={onDelete}
              >
                <TrashCanIcon className="w-3 h-3" />
              </Button>
            </div>
          </div>
        </CardBody>
      </Card>
    );
  }

  function SetsList({ sets }: { sets: WorkoutSet[] }) {
    return (
      <div className="flex flex-col gap-4">
        {sets.map((set, index) => (
          <SetCard key={set.id} num={index + 1} set={set} />
        ))}
      </div>
    );
  }

  function ContentCard() {
    async function onAddSet() {
      await authApi.v1
        .routineServiceAddSetToExerciseInstance(id, exerciseInstanceId, {
          reps: exerciseInstanceDetails.sets![
            exerciseInstanceDetails.sets!.length - 1
          ].reps,
          setType: WorkoutSetType.SET_TYPE_REPS,
        })
        .then((response) => {
          console.log(response.data);
          fetchExerciseInstanceDetails();
        })
        .catch((error) => {
          console.log(error);
          toast.error("Не удалось добавить подход");
        });
    }

    function AddSetBlock() {
      return (
        <Card>
          <CardBody className="flex flex-col gap-4">
            <h2 className="text-md font-semibold">
              скибиди туалет я еще не приудмал че суда написать
            </h2>
          </CardBody>
        </Card>
      );
    }

    return (
      <div className="p-4 h-full overflow-y-auto flex flex-col gap-4">
        <AddSetBlock />
        {exerciseInstanceDetails.sets!.length > 0 && (
          <SetsList sets={exerciseInstanceDetails.sets!} />
        )}
        <div className="flex justify-center">
          <Button fullWidth color="primary" onPress={onAddSet}>
            <PlusIcon className="w-4 h-4" />
            Добавить подход
          </Button>
        </div>
      </div>
    );
  }

  async function onDelete() {
    try {
      await authApi.v1
        .routineServiceRemoveExerciseInstanceFromRoutine(id, exerciseInstanceId)
        .then((response) => {
          console.log(response.data);
          router.push(`/routines/${id}`);
        })
        .catch((error) => {
          console.log(error);
          if (error === errUnauthorized || error.response?.status === 401) {
            router.push("/auth/login");

            return;
          }
          throw error;
        });
    } catch (error) {
      console.log(error);
      toast.error("Не удалось удалить инстанс упражнения");
    }
  }

  return (
    <div className="py-4 flex flex-col h-full">
      <PageHeader
        enableBackButton
        title={exerciseInstanceDetails.exercise?.name!}
      >
        <DropdownItem
          key="delete"
          className="text-danger"
          color="danger"
          onPress={onDelete}
        >
          Удалить
        </DropdownItem>
      </PageHeader>
      <ContentCard />
    </div>
  );
}
