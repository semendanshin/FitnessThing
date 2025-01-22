"use client";

import { authApi, errUnauthorized } from "@/api/api";
import { WorkoutExerciseLogDetails, WorkoutSetLog } from "@/api/api.generated";
import { Loading } from "@/components/loading";
import { PlusIcon, TrashCanIcon } from "@/config/icons";
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
import { PageHeader } from "@/components/page-header";
import { Divider } from "@nextui-org/divider";

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
                if (
                    error === errUnauthorized ||
                    error.response?.status === 401
                ) {
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
                setExerciseLogHistory(response.data.exerciseLogs!);
            })
            .catch((error) => {
                console.log(error);
                if (
                    error === errUnauthorized ||
                    error.response?.status === 401
                ) {
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
                <h2 className="text-lg text-red-500">
                    Ошибка при загрузке данных
                </h2>
                <p>Проверьте соединение с сервером или обновите страницу.</p>
            </div>
        );
    }

    function SetLogCard({
        setLog,
        setNum,
        enableDelete,
        onDelete,
    }: {
        setLog: WorkoutSetLog;
        setNum: number;
        enableDelete?: boolean;
        onDelete?: () => void;
    }) {
        return (
            <Card className="flex flex-row items-center justify-between p-2 w-full">
                <div className="flex flex-row w-full gap-2 px-2">
                    <div className="text-sm font-semibold w-4">
                        {setNum + 1}
                    </div>
                    <div className="text-sm font-semibold w-fit">
                        {setLog?.weight} кг
                    </div>
                    <div className="text-sm font-semibold w-fit">x</div>
                    <div className="text-sm font-semibold w-fit">
                        {setLog?.reps} раз
                    </div>
                </div>
                {enableDelete && (
                    <div className="flex flex-col">
                        <Button
                            color="danger"
                            size="sm"
                            onPress={onDelete}
                            isIconOnly
                            className="h-fit w-fit min-w-fit p-2"
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

    function AddSetLogModal({
        isOpen,
        onClose,
    }: {
        isOpen: boolean;
        onClose: () => void;
    }) {
        const [weight, setWeight] = useState<number>(0);
        const [reps, setReps] = useState<number>(0);

        const [errors, setErrors] = useState<{
            weight?: string;
            reps?: string;
        }>({});

        async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
            event.preventDefault();
            setIsLoading(true);
            try {
                await authApi.v1.workoutServiceLogSet(id, exerciseLogId, {
                    weight: weight!,
                    reps: reps!,
                });
                onClose();
                await fetchData();
            } catch (error) {
                console.log(error);
                toast.error("Failed to add exercises to workout");
            } finally {
                setIsLoading(false);
            }
        }

        return (
            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalContent>
                    {(onClose) => (
                        <div className="flex flex-col py-4 mb-4">
                            <ModalHeader className="p-0 px-4">
                                Добавить сет
                            </ModalHeader>
                            <Form
                                validationBehavior="native"
                                validationErrors={errors}
                                onSubmit={handleSubmit}
                                className="flex flex-col p-0 px-2"
                            >
                                <ModalBody className="flex flex-row gap-2 px-2 w-full">
                                    <Input
                                        type="number"
                                        label="Вес"
                                        placeholder="Вес"
                                        onValueChange={(value) =>
                                            setWeight(Number(value))
                                        }
                                        isRequired
                                        autoFocus
                                    />
                                    <Input
                                        type="number"
                                        label="Повторы"
                                        placeholder="Повторы"
                                        onValueChange={(value) =>
                                            setReps(Number(value))
                                        }
                                        isRequired
                                    />
                                </ModalBody>
                                <ModalFooter className="flex flex-col gap-2 w-full justify-around px-2 py-0">
                                    <Button
                                        color="success"
                                        type="submit"
                                        className="w-full"
                                    >
                                        Добавить
                                    </Button>
                                    <Button
                                        color="danger"
                                        onPress={onClose}
                                        className="w-full"
                                    >
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
        const [weight, setWeight] = useState<number>(
            exerciseLogHistory[0]?.setLogs![0]?.weight!
        );
        const [reps, setReps] = useState<number>(
            exerciseLogHistory[0]?.setLogs![0]?.reps!
        );

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
                        onPress={() => {
                            if (value > 0 && isSubtract) {
                                setValue(value - 1);
                                return;
                            }
                            if (!isSubtract) {
                                setValue(value + 1);
                            }
                        }}
                        className="min-w-fit w-fit p-3"
                        isIconOnly
                    >
                        {isSubtract ? "-" : "+"}
                    </Button>
                </div>
            );
        }

        async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
            event.preventDefault();
            try {
                await authApi.v1.workoutServiceLogSet(id, exerciseLogId, {
                    weight: weight!,
                    reps: reps!,
                });
                await fetchData();
            } catch (error) {
                console.log(error);
                toast.error("Failed to add exercises to workout");
            }
        }

        async function onDeleteSet(setId: string) {
            try {
                await authApi.v1.workoutServiceDeleteSetLog(
                    id,
                    exerciseLogId,
                    setId
                );
                await fetchData();
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
                            <label>Вес</label>
                            <div className="flex flex-row gap-2 items-center">
                                <IncrementButtons
                                    value={weight}
                                    setValue={setWeight}
                                    isSubtract
                                />
                                <Input
                                    type="number"
                                    placeholder="10"
                                    className="p-0 w-full h-full"
                                    isRequired
                                    value={String(weight)}
                                    onValueChange={(value) =>
                                        setWeight(Number(value))
                                    }
                                />
                                <IncrementButtons
                                    value={weight}
                                    setValue={setWeight}
                                />
                            </div>
                        </div>
                        <div className="flex flex-col gap-1 w-1/2">
                            <label>Повторы</label>
                            <div className="flex flex-row gap-2 items-center">
                                <IncrementButtons
                                    value={reps}
                                    setValue={setReps}
                                    isSubtract
                                />
                                <Input
                                    type="number"
                                    placeholder="10"
                                    className="p-0 w-full h-full"
                                    isRequired
                                    value={String(reps)}
                                    onValueChange={(value) =>
                                        setReps(Number(value))
                                    }
                                />
                                <IncrementButtons
                                    value={reps}
                                    setValue={setReps}
                                />
                            </div>
                        </div>
                    </div>
                    <Button
                        color="primary"
                        className="w-full"
                        type="submit"
                        size="sm"
                    >
                        Добавить
                    </Button>
                </Form>
                <Divider />
                <div className="flex flex-col gap-2">
                    {exerciseLogDetails.setLogs?.map((setLog, index) => (
                        <SetLogCard
                            key={index}
                            setLog={setLog}
                            setNum={index}
                            enableDelete
                            onDelete={() => onDeleteSet(setLog.id!)}
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
                                <div
                                    key={index}
                                    className="flex flex-col gap-2"
                                >
                                    <label className="text-lg font-bold">
                                        {new Date(
                                            exerciseLog.exerciseLog?.createdAt!
                                        ).toLocaleDateString("ru-RU", {
                                            weekday: "long",
                                            day: "numeric",
                                            month: "long",
                                        })}
                                    </label>
                                    <div className="flex flex-col gap-2">
                                        {exerciseLog.setLogs?.map(
                                            (setLog, index) => (
                                                <SetLogCard
                                                    key={index}
                                                    setLog={setLog}
                                                    setNum={index}
                                                />
                                            )
                                        )}
                                    </div>
                                </div>
                            )
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
                        title={exerciseLogDetails.exercise?.name!}
                        enableBackButton
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
            <AddSetLogModal isOpen={isOpen} onClose={onClose} />
        </>
    );
}
