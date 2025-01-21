"use client";

import { authApi, errUnauthorized } from "@/api/api";
import {
    WorkoutExerciseLogDetails,
    WorkoutGetWorkoutResponse,
} from "@/api/api.generated";
import { Loading } from "@/components/loading";
import { PageHeader } from "@/components/page-header";
import { ModalSelectExercise } from "@/components/pick-exercises-modal";
import { LeftArrowIcon, PlusIcon, RightArrowIcon } from "@/config/icons";
import { Button } from "@nextui-org/button";
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { Link } from "@nextui-org/link";
import { useDisclosure } from "@nextui-org/modal";
import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";
import { toast } from "react-toastify";

export default function RoutineDetailsPage({
    params,
}: {
    params: Promise<{ id: string }>;
}) {
    const [isLoading, setIsLoading] = useState(true);
    const [isError, setIsError] = useState(false);

    const [workoutDetails, setWorkoutDetails] =
        useState<WorkoutGetWorkoutResponse>({});

    const { id } = use(params);

    const router = useRouter();

    const { isOpen, onOpen, onClose } = useDisclosure();

    async function fetchWorkoutDetails() {
        await authApi.v1
            .workoutServiceGetWorkout(id)
            .then((response) => {
                console.log(response.data);
                setWorkoutDetails(response.data!);
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
            await fetchWorkoutDetails();
        } catch (error) {
            toast.error("Failed to fetch workout details");
            setIsError(true);
        } finally {
            setIsLoading(false);
        }
    }

    async function finishWorkout() {
        setIsLoading(true);
        try {
            await authApi.v1
                .workoutServiceCompleteWorkout(id, {})
                .then((response) => {
                    console.log(response.data);
                    router.push(`/workouts/${id}/results`);
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
        } catch (error) {
            toast.error("Не удалось завершить тренировку");
        } finally {
            setIsLoading(false);
        }
    }

    async function addExercisesToWorkout(exerciseIds: string[]) {
        try {
            await Promise.all(
                exerciseIds.map((exerciseId) =>
                    authApi.v1
                        .workoutServiceLogExercise(id, {
                            exerciseId,
                        })
                        .then((response) => {
                            console.log(response.data);
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
                        })
                )
            );
            onClose();
            await fetchData();
        } catch (error) {
            console.log(error);
            toast.error("Failed to add exercises to workout");
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
                <h2 className="text-lg text-red-500">
                    Ошибка при загрузке данных
                </h2>
                <p>Проверьте соединение с сервером или обновите страницу.</p>
            </div>
        );
    }

    function ExerciseLogCard({
        exerciseLogDetails,
    }: {
        exerciseLogDetails: WorkoutExerciseLogDetails;
    }) {
        return (
            <Card
                className="flex flex-row items-center justify-between p-2 cursor-pointer"
                shadow="sm"
                fullWidth
                as={Link}
                href={`/workouts/${id}/exerciseLogs/${exerciseLogDetails.exerciseLog?.id}`}
            >
                <div className="flex flex-col items-start justify-between p-2 w-full">
                    <CardHeader className="p-0">
                        <label className="text-m font-bold">
                            {exerciseLogDetails.exercise?.name}
                        </label>
                    </CardHeader>
                    {exerciseLogDetails.setLogs!.length > 0 && (
                        <CardBody className="w-full p-0 py-2">
                            <div className="grid grid-cols-3 gap-2 w-fit">
                                {exerciseLogDetails.setLogs?.map(
                                    (setLog, index) => (
                                        <div key={index} className="contents">
                                            <div className="text-sm font-semibold">
                                                {index + 1}
                                            </div>
                                            <div className="text-sm font-semibold">
                                                {setLog?.weight} кг
                                            </div>
                                            <div className="text-sm font-semibold">
                                                {setLog?.reps} раз
                                            </div>
                                        </div>
                                    )
                                )}
                            </div>
                        </CardBody>
                    )}
                </div>
                <div className="flex flex-col items-center justify-between p-2">
                    <RightArrowIcon className="w-4 h-4" fill="currentColor" />
                </div>
            </Card>
        );
    }

    return (
        <>
            <div className="py-4 flex-grow max-w-full">
                <div className="h-full max-h-full overflow-y-auto gap-4 flex flex-col">
                    <PageHeader title={"Тренировка"} enableBackButton={true} />
                    <section className="flex flex-col gap-4 px-4">
                        <div className="flex flex-col gap-2">
                            {workoutDetails.exerciseLogs?.map(
                                (exerciseLogDetails, index) => (
                                    <ExerciseLogCard
                                        key={index}
                                        exerciseLogDetails={exerciseLogDetails}
                                    />
                                )
                            )}
                            <Card className="p-2">
                                <Button
                                    className="w-full"
                                    onPress={() => {
                                        onOpen();
                                    }}
                                >
                                    <PlusIcon className="w-5 h-5" />
                                    <span>Добавить упражнение</span>
                                </Button>
                            </Card>
                        </div>
                    </section>
                    <section className="w-full bottom-0 px-4">
                        <Button
                            color="secondary"
                            className="w-full"
                            onPress={async () => {
                                await finishWorkout();
                            }}
                        >
                            Завершить тренировку
                        </Button>
                    </section>
                </div>
            </div>
            <ModalSelectExercise
                isOpen={isOpen}
                onClose={onClose}
                excludeExerciseIds={workoutDetails.exerciseLogs!.map(
                    (exerciseLog) => exerciseLog.exerciseLog!.exerciseId!
                )}
                onSubmit={addExercisesToWorkout}
            />
        </>
    );
}
