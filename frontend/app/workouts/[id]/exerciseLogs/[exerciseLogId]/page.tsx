"use client";

import { Button } from "@nextui-org/button";
import { Card, CardBody } from "@nextui-org/card";
import { Form } from "@nextui-org/form";
import { Input } from "@nextui-org/react";
import { Tabs, Tab } from "@nextui-org/tabs";
import { DropdownItem } from "@nextui-org/dropdown";
import {
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/modal";
import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";
import { toast } from "react-toastify";
import { Divider } from "@nextui-org/divider";

import { PageHeader } from "@/components/page-header";
import { TrashCanIcon } from "@/config/icons";
import { Loading } from "@/components/loading";
import { WorkoutExerciseLogDetails, WorkoutSetLog } from "@/api/api.generated";
import { authApi, errUnauthorized } from "@/api/api";

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string; exerciseLogId: string }>;
}) {
  const { exerciseLogId, id } = use(params);

  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const [exerciseLogDetails, setExerciseLogDetails] =
    useState<WorkoutExerciseLogDetails>({});
  const [exerciseLogHistory, setExerciseLogHistory] = useState<
    WorkoutExerciseLogDetails[]
  >([]);
  const [exerciseLogForUpdate, setExerciseLogForUpdate] =
    useState<WorkoutSetLog>({});

  const { isOpen, onOpen, onClose } = useDisclosure();

  const router = useRouter();

  async function fetchExerciseLogDetails() {
    await authApi.v1
      .workoutServiceGetExerciseLogDetails(id, exerciseLogId)
      .then((response) => {
        console.log(response.data);
        setExerciseLogDetails(response.data.exerciseLogDetails!);
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

  async function fetchExerciseLogHistory(exerciseId: string) {
    await authApi.v1
      .exerciseServiceGetExerciseHistory(exerciseId)
      .then((response) => {
        console.log(response.data);
        setExerciseLogHistory(
          response.data.exerciseLogs!.filter(
            (log) => log.exerciseLog?.workoutId != id,
          ),
        );
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
      await fetchExerciseLogDetails();
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

  useEffect(() => {
    if (exerciseLogDetails.exercise?.id) {
      fetchExerciseLogHistory(exerciseLogDetails.exercise.id);
    }
  }, [exerciseLogDetails.exercise?.id]);

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

  function SetLogCard({
    setLog,
    setNum,
    enableDelete,
    onDelete,
    onPress,
  }: {
    setLog: WorkoutSetLog;
    setNum: number;
    enableDelete?: boolean;
    onDelete?: () => Promise<void>;
    onPress?: () => void;
  }) {
    const [isLoading, setIsLoading] = useState(false);

    return (
      <Card
        className="flex flex-row items-center justify-between p-2 w-full"
        isPressable={!!onPress}
        onPress={onPress}
      >
        <div className="flex flex-row w-full gap-2 px-2">
          <div className="text-sm font-semibold w-4">{setNum + 1}</div>
          <div className="text-sm font-semibold w-fit">{setLog?.weight} кг</div>
          <div className="text-sm font-semibold w-fit">x</div>
          <div className="text-sm font-semibold w-fit">{setLog?.reps} раз</div>
        </div>
        {enableDelete && (
          <div className="flex flex-col">
            <Button
              isIconOnly
              className="h-fit w-fit min-w-fit p-2 w-9 h-6"
              color="danger"
              isLoading={isLoading}
              size="sm"
              onPress={async () => {
                setIsLoading(true);
                await onDelete!();
                setIsLoading(false);
              }}
            >
              <TrashCanIcon className="w-3 h-3" />
            </Button>
          </div>
        )}
      </Card>
    );
  }

  async function onDelete() {
    try {
      await authApi.v1.workoutServiceDeleteExerciseLog(id, exerciseLogId);
      router.back();
    } catch (error) {
      console.log(error);
      toast.error("Failed to delete exercise log");
    }
  }

  function IncrementButtons({
    value,
    setValue,
    isSubtract,
  }: {
    value: number;
    setValue: (value: number) => void;
    isSubtract?: boolean;
  }) {
    return (
      <div className="flex flex-col justify-around p-0">
        <Button
          isIconOnly
          className="min-w-fit w-fit p-3"
          onPress={() => {
            if (value > 0 && isSubtract) {
              setValue(value - 1);

              return;
            }
            if (!isSubtract) {
              setValue(value + 1);
            }
          }}
        >
          {isSubtract ? "-" : "+"}
        </Button>
      </div>
    );
  }

  function InputWithIncrement({
    value,
    setValue,
    label,
    placeholder,
    type,
  }: {
    value: number;
    setValue: (value: number) => void;
    label: string;
    placeholder: string;
    type: string;
  }) {
    return (
      <>
        <p>{label}</p>
        <div className="flex flex-row gap-2 items-center">
          <IncrementButtons isSubtract setValue={setValue} value={value} />
          <Input
            isRequired
            className="p-0 w-full h-full"
            placeholder={placeholder}
            type={type}
            value={String(value)}
            onValueChange={(value) => setValue(Number(value))}
          />
          <IncrementButtons setValue={setValue} value={value} />
        </div>
      </>
    );
  }

  function UpdateSetLogModal({
    isOpen,
    onClose,
    setLog,
  }: {
    isOpen: boolean;
    onClose: () => void;
    setLog: WorkoutSetLog;
  }) {
    const [weight, setWeight] = useState<number>(setLog.weight!);
    const [reps, setReps] = useState<number>(setLog.reps!);

    const [errors, setErrors] = useState<{
      weight?: string;
      reps?: string;
    }>({});

    async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
      event.preventDefault();
      try {
        await authApi.v1.workoutServiceUpdateSetLog(
          id,
          exerciseLogId,
          setLog.id!,
          {
            weight: weight!,
            reps: reps!,
          },
        );
        fetchExerciseLogDetails();
        onClose();
      } catch (error) {
        console.log(error);
        toast.error("Failed to update set log");
      } finally {
        setIsLoading(false);
      }
    }

    return (
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          {(onClose) => (
            <div className="flex flex-col py-4 mb-4">
              <ModalHeader className="p-0 px-4">Изменить сет</ModalHeader>
              <Form
                className="flex flex-col p-0 px-2"
                validationBehavior="native"
                validationErrors={errors}
                onSubmit={handleSubmit}
              >
                <ModalBody className="flex flex-row gap-2 px-2 w-full">
                  <div className="flex flex-col gap-1 w-1/2">
                    <InputWithIncrement
                      label="Вес"
                      placeholder="10"
                      setValue={setWeight}
                      type="number"
                      value={weight}
                    />
                  </div>
                  <div className="flex flex-col gap-1 w-1/2">
                    <InputWithIncrement
                      label="Повторы"
                      placeholder="10"
                      setValue={setReps}
                      type="number"
                      value={reps}
                    />
                  </div>
                </ModalBody>
                <ModalFooter className="flex flex-col gap-2 w-full justify-around px-2 py-0">
                  <Button className="w-full" color="success" type="submit">
                    Изменить
                  </Button>
                  <Button className="w-full" color="danger" onPress={onClose}>
                    Отмена
                  </Button>
                </ModalFooter>
              </Form>
            </div>
          )}
        </ModalContent>
      </Modal>
    );
  }

  function TodayContent() {
    const [weight, setWeight] = useState<number>(0);
    const [reps, setReps] = useState<number>(0);

    useEffect(() => {
      if (exerciseLogDetails.setLogs?.length) {
        console.log("из текущей тренировки");
        const lastIndex = exerciseLogDetails.setLogs.length - 1;

        setWeight(exerciseLogDetails.setLogs[lastIndex]?.weight!);
        setReps(exerciseLogDetails.setLogs[lastIndex]?.reps!);
      } else if (exerciseLogHistory.length) {
        console.log("из истории");
        const lastIndex = exerciseLogHistory.length - 1;

        setWeight(exerciseLogHistory[0]!.setLogs![lastIndex]?.weight!);
        setReps(exerciseLogHistory[0]!.setLogs![lastIndex]?.reps!);
      }
    }, [exerciseLogDetails, exerciseLogHistory]);

    async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
      event.preventDefault();
      try {
        await authApi.v1.workoutServiceLogSet(id, exerciseLogId, {
          weight: weight!,
          reps: reps!,
        });
        await fetchExerciseLogDetails();
      } catch (error) {
        console.log(error);
        toast.error("Failed to add exercises to workout");
      }
    }

    async function onDeleteSet(setId: string) {
      try {
        await authApi.v1.workoutServiceDeleteSetLog(id, exerciseLogId, setId);
        await fetchExerciseLogDetails();
      } catch (error) {
        console.log(error);
        toast.error("Failed to delete set");
      }
    }

    return (
      <div className="flex flex-col gap-4">
        <Form className="flex flex-col gap-3" onSubmit={handleSubmit}>
          <div className="flex flex-row justify-around gap-4">
            <div className="flex flex-col gap-1 w-1/2">
              <InputWithIncrement
                label="Вес"
                placeholder="10"
                setValue={setWeight}
                type="number"
                value={weight}
              />
            </div>
            <div className="flex flex-col gap-1 w-1/2">
              <InputWithIncrement
                label="Повторы"
                placeholder="10"
                setValue={setReps}
                type="number"
                value={reps}
              />
            </div>
          </div>
          <Button className="w-full" color="primary" size="sm" type="submit">
            Добавить
          </Button>
        </Form>
        <Divider />
        <div className="flex flex-col gap-2">
          {exerciseLogDetails.setLogs?.map((setLog, index) => (
            <SetLogCard
              key={index}
              enableDelete
              setLog={setLog}
              setNum={index}
              onDelete={() => onDeleteSet(setLog.id!)}
              onPress={() => {
                setExerciseLogForUpdate(setLog);
                onOpen();
              }}
            />
          ))}
        </div>
      </div>
    );
  }

  function HistoryContent() {
    return (
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          {exerciseLogHistory.map(
            (exerciseLog, index) =>
              exerciseLog.setLogs!.length > 0 &&
              exerciseLog.exerciseLog?.workoutId != id && (
                <div key={index} className="flex flex-col gap-2">
                  <p className="text-lg font-bold">
                    {new Date(
                      exerciseLog.exerciseLog?.createdAt!,
                    ).toLocaleDateString("ru-RU", {
                      weekday: "long",
                      day: "numeric",
                      month: "long",
                    })}
                  </p>
                  <div className="flex flex-col gap-2">
                    {exerciseLog.setLogs?.map((setLog, index) => (
                      <SetLogCard key={index} setLog={setLog} setNum={index} />
                    ))}
                  </div>
                </div>
              ),
          )}
        </div>
      </div>
    );
  }

  function TabContent({ children }: { children: JSX.Element }) {
    return (
      <Card>
        <CardBody>{children}</CardBody>
      </Card>
    );
  }

  return (
    <>
      <div className="py-4 flex-grow max-w-full">
        <div className="h-full max-h-full overflow-y-auto gap-4 flex flex-col">
          <PageHeader
            enableBackButton
            title={exerciseLogDetails.exercise?.name!}
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
          <section className="flex flex-col flex-grow gap-4 px-4 justify-start">
            <Tabs aria-label="Options">
              <Tab key="today" title="Сегодня">
                <TabContent>
                  <TodayContent />
                </TabContent>
              </Tab>
              <Tab key="history" title="История">
                <TabContent>
                  <HistoryContent />
                </TabContent>
              </Tab>
            </Tabs>
          </section>
        </div>
      </div>
      {/* <AddSetLogModal isOpen={isOpen} onClose={onClose} /> */}
      <UpdateSetLogModal
        isOpen={isOpen}
        setLog={exerciseLogForUpdate}
        onClose={onClose}
      />
    </>
  );
}
