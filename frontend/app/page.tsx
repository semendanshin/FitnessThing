/* eslint-disable react/no-unknown-property */
"use client";
import { Button } from "@nextui-org/button";
import { Card, CardBody, CardFooter, CardHeader } from "@nextui-org/card";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";

import { WorkoutWorkout } from "@/api/api.generated";
import { AnimationProcessor } from "@/components/animated-background";
import { BoltIcon, ChevronRightIcon, PlayIcon } from "@/config/icons";
import { Loading } from "@/components/loading";
import { authApi, errUnauthorized } from "@/api/api";

export default function Home() {
  const [user, setUser] = useState<any>({});
  const [routines, setRoutines] = useState<any[]>([]);
  const [activeWorkouts, setActiveWorkouts] = useState<WorkoutWorkout[]>([]);

  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const router = useRouter();

  async function fetchUser() {
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

  async function fetRoutines() {
    await authApi.v1
      .routineServiceGetRoutines()
      .then((response) => {
        console.log(response.data);
        setRoutines(response.data.routines!);
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

  async function fetchActiveWorkouts() {
    await authApi.v1
      .workoutServiceGetActiveWorkouts()
      .then((response) => {
        console.log(response.data);
        setActiveWorkouts(response.data.workouts!);
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
      await Promise.all([fetchUser(), fetRoutines(), fetchActiveWorkouts()]);
      setIsError(false);
    } catch (error) {
      console.log(error);
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
  }

  async function startWorkout(routineId: string | undefined) {
    if (activeWorkouts.length > 0) {
      toast.error("Сначала завершите активную тренировку");

      return;
    }
    await authApi.v1
      .workoutServiceStartWorkout({
        routineId: routineId,
      })
      .then((response) => {
        console.log(response.data);
        router.push(`/workouts/${response.data.workout?.id}`);
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

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    const canvas = document.getElementById("home-bg") as HTMLCanvasElement;

    if (!canvas) return;

    const animationProcessor = new AnimationProcessor(canvas, 400, 1);

    const updateAnimationDimensions = () => {
      const container = canvas.parentElement?.parentElement;

      if (!container) return;

      const width = container.clientWidth;
      const height = container.clientHeight;

      // Update center to be in the middle of the container
      animationProcessor.updateCenter(width / 2, height / 2);

      // Update radius to be 60% of the smallest dimension
      const radius = Math.min(width, height) * 0.6;

      animationProcessor.updateRadius(radius);
    };

    // Initial setup
    updateAnimationDimensions();
    animationProcessor.start();

    // Add resize listener
    window.addEventListener("resize", updateAnimationDimensions);

    return () => {
      animationProcessor.stop();
      window.removeEventListener("resize", updateAnimationDimensions);
    };
  }, [isLoading]);

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
    <div className="flex flex-col flex-grow max-w-full justify-start">
      <section
        className="flex flex-col flex-grow max-w-full h-[65vh]"
        id="home"
      >
        <div className="flex flex-col flex-grow absolute w-full h-[75vh] z-0 fade-bottom opacity-80">
          <canvas id="home-bg" />
        </div>
        <h1 className="text-2xl font-bold p-4 absolute">
          Привет{user.firstName ? `, ${user.firstName}` : ""}!
        </h1>
        <div className="flex-grow flex flex-col justify-center items-center relative z-10 drop-shadow-lg">
          <Button
            disableRipple
            className="flex items-center bg-transparent h-fit"
            size="lg"
            onPress={async () => {
              await startWorkout(undefined);
            }}
          >
            <BoltIcon className="w-7 h-7" fill="currentColor" />
            <span className="text-2xl font-bold">Стать лучше</span>
          </Button>
          <Button
            className="flex items-center text-white-500 bg-transparent underline p-0"
            size="sm"
            onPress={async () => {
              await startWorkout(undefined);
            }}
          >
            Пустая тренировка
          </Button>
        </div>
      </section>
      {activeWorkouts.length > 0 && (
        <section className="flex flex-col flex-grow max-w-full h-fit relative">
          <h4 className="text-xl font-bold px-4">Активные тренировки</h4>
          <div className="flex flex-col p-4 max-w-full overflow-y-auto">
            {activeWorkouts.map((workout) => (
              <Card key={workout.id} className="w-full">
                <CardHeader>
                  <h3 className="text-lg font-bold">
                    {"Тренировка от "}
                    {new Date(workout.createdAt!).toLocaleString("ru-RU", {
                      weekday: "long",
                      day: "numeric",
                      month: "long",
                    })}
                  </h3>
                </CardHeader>
                <CardFooter>
                  <Button
                    className="flex items-center px-2 w-full"
                    color="primary"
                    size="sm"
                    onPress={async () => {
                      router.push(`/workouts/${workout.id}`);
                    }}
                  >
                    <PlayIcon className="w-3 h-3" fill="currentColor" />
                    <span className="text-sm font-bold">Продолжить</span>
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        </section>
      )}
      <section className="flex flex-col flex-grow max-w-full h-fit relative max-w-full">
        {/* Список шаблонов. Горизонтальный скролл с квадратными карточками */}
        <Link className="flex items-center px-4 gap-1" href="/routines">
          <h4 className="text-xl font-bold">Шаблоны</h4>
          <ChevronRightIcon className="w-4 h-4" fill="currentColor" />
        </Link>
        <div className="flex flex-col p-4 max-w-full overflow-x-auto">
          <div className="flex flex-row gap-4 justify-start w-fit">
            {routines.map((routine) => (
              <Card
                key={routine.id}
                as={Link}
                className="w-52 h-52"
                href={`/routines/${routine.id}`}
              >
                <CardHeader>
                  <h3 className="text-lg font-bold">{routine.name}</h3>
                </CardHeader>
                <CardBody>
                  <p>{routine.description}</p>
                </CardBody>
                <CardFooter>
                  <Button
                    className="flex items-center px-2 w-full"
                    color="primary"
                    size="sm"
                    onPress={async () => {
                      await startWorkout(routine.id);
                    }}
                  >
                    <PlayIcon className="w-3 h-3" fill="currentColor" />
                    <span className="text-sm font-bold">Начать</span>
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        </div>
      </section>
      <style jsx>{`
        .fade-bottom::after {
          content: "";
          position: absolute;
          bottom: 0;
          left: 0;
          width: 100%;
          height: 50px;
          background: linear-gradient(
            to top,
            hsl(var(--nextui-background)),
            transparent 100%
          );
          pointer-events: none;
        }
      `}</style>
    </div>
  );
}
