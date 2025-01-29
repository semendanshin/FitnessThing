import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { Checkbox } from "@nextui-org/checkbox";
import { Input } from "@nextui-org/input";
import { Select, SelectItem } from "@nextui-org/select";
import { Skeleton } from "@nextui-org/skeleton";
import { Modal, ModalContent, ModalHeader } from "@nextui-org/modal";
import { ScrollShadow } from "@nextui-org/scroll-shadow";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { Button } from "@nextui-org/button";

import InfiniteScroll, { useInfiniteScroll } from "./infinite-scroll";

import { authApi, errUnauthorized } from "@/api/api";

function SkeletonExerciseCard() {
  return (
    <Card fullWidth className="flex flex-row flex-grow p-2 gap-4">
      <div className="flex flex-col items-start justify-between p-2 w-full">
        <CardHeader className="p-0">
          <Skeleton className="w-4/5 rounded-lg">
            <div className="h-4 rounded-lg bg-default-200" />
          </Skeleton>
        </CardHeader>
        <CardBody className="w-full px-0 py-2">
          <div className="flex flex-row gap-2 w-full">
            <Skeleton className="w-2/3 rounded-lg">
              <div className="h-3 rounded-lg bg-default-200" />
            </Skeleton>
          </div>
        </CardBody>
      </div>
      <div className="flex items-center justify-end p-2 w-full">
        <Skeleton className="w-2/3 rounded-lg">
          <div className="h-7 rounded-lg bg-default-200" />
        </Skeleton>
      </div>
    </Card>
  );
}

function ExerciseCard({
  exercise,
  isSelected,
  onSelectedChange,
}: {
  exercise: any;
  isSelected: boolean;
  onSelectedChange: () => void;
}) {
  return (
    <Card
      key={exercise.id}
      fullWidth
      className="flex flex-row flex-grow p-2 gap-4 justify-between cursor-pointer max-h-fit"
      shadow="sm"
    >
      <div className="flex flex-col items-start justify-between p-2">
        <CardHeader className="p-0">
          <p className="text-m font-bold">{exercise.name}</p>
        </CardHeader>
        <CardBody className="p-0">
          {exercise.description ? (
            <div className="py-2">
              <p className="text-xs text-gray-400/80">{exercise.description}</p>
            </div>
          ) : null}
          {exercise.targetMuscleGroups ? (
            <div className="flex flex-row gap-2">
              <p className="text-xs text-gray-400">
                {exercise.targetMuscleGroups.join(", ")}
              </p>
            </div>
          ) : null}
        </CardBody>
      </div>
      <div className="flex flex-col items-center justify-center p-2">
        <Checkbox checked={isSelected} onChange={onSelectedChange} />
      </div>
    </Card>
  );
}

export function ModalSelectExercise({
  isOpen,
  onClose,
  excludeExerciseIds,
  onSubmit,
}: {
  isOpen: boolean;
  onClose: () => void;
  excludeExerciseIds: string[];
  onSubmit: (selectedExercisesIds: string[]) => void;
}) {
  const [serachQuery, setSearchQuery] = useState("");
  const [muscleGroup, setMuscleGroup] = useState("");

  const { hasMore, setHasMore } = useInfiniteScroll();

  const [exercises, setExercises] = useState<any[]>([]);
  const [muscleGroups, setMuscleGroups] = useState<any[]>([]);
  const [selectedExercisesIds, setSelectedExercisesIds] = useState<string[]>(
    [],
  );

  const [isLoading, setIsLoading] = useState(true);

  const router = useRouter();

  function toggleExerciseSelection(id: string) {
    setSelectedExercisesIds((prev) => {
      const index = prev.indexOf(id);

      if (index === -1) {
        return [...prev, id];
      }

      return [...prev.slice(0, index), ...prev.slice(index + 1)];
    });
  }

  async function fetchExercises() {
    if (!isOpen) {
      return;
    }
    setIsLoading(true);
    await authApi.v1
      .exerciseServiceGetExercises(
        {
          muscleGroupIds: muscleGroup ? [muscleGroup] : undefined,
          excludeExerciseIds: excludeExerciseIds,
        },
        {
          paramsSerializer: {
            indexes: null,
          },
        },
      )
      .then((response) => {
        console.log(response.data);
        setExercises(response.data.exercises!);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при загрузке упражнений");
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  async function fetchMuscleGroups() {
    if (!isOpen) {
      return;
    }
    await authApi.v1
      .exerciseServiceGetMuscleGroups()
      .then((response) => {
        console.log(response.data);
        setMuscleGroups(response.data.muscleGroups!);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при загрузке групп мышц");
      });
  }

  const fetchMore = async () => {
    setHasMore(false);
    // await authApi.v1
    //     .exerciseServiceGetExercises(
    //         {
    //             muscleGroupIds: muscleGroup ? [muscleGroup] : undefined,
    //             excludeExerciseIds: excludeExerciseIds,
    //             // offset: offset,
    //             // limit: limit,
    //         },
    //         {
    //             paramsSerializer: {
    //                 indexes: null,
    //             },
    //         }
    //     )
    //     .then((response) => {
    //         console.log(response.data);
    //         setExercises((prev) => [...prev, ...response.data.exercises!]);
    //         setHasMore(response.data.exercises!.length === limit);
    //     })
    //     .catch((error) => {
    //         console.log(error);
    //         if (
    //             error === errUnauthorized ||
    //             error.response?.status === 401
    //         ) {
    //             router.push("/auth/login");
    //             return;
    //         }
    //         toast.error("Ошибка при загрузке упражнений");
    //     });
  };

  useEffect(() => {
    fetchMuscleGroups();
  }, [isOpen]);

  useEffect(() => {
    fetchExercises();
  }, [muscleGroup, isOpen]);

  return (
    <Modal
      className="overflow-y-auto h-full p-2 w-full h-[85vh]"
      isOpen={isOpen}
      placement="center"
      scrollBehavior="inside"
      size="xs"
      onClose={onClose}
    >
      <ModalContent className="h-full">
        {(onClose) => (
          <div className="flex flex-col h-full">
            <ModalHeader className="p-2">Выберите упражнение</ModalHeader>

            <ScrollShadow className="flex-grow">
              <div className="flex flex-col gap-4 p-2 flex-grow">
                <Input
                  placeholder="Поиск"
                  value={serachQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
                <Select
                  aria-label="Выберите группу мышц"
                  placeholder="Выберите упражнение"
                  onChange={(e) => {
                    setMuscleGroup(e.target.value);
                  }}
                >
                  {muscleGroups.map((muscleGroup) => (
                    <SelectItem key={muscleGroup.id} value={muscleGroup.id}>
                      {muscleGroup.name}
                    </SelectItem>
                  ))}
                </Select>
                <InfiniteScroll
                  showError
                  showLoading
                  className="flex flex-col gap-2"
                  fetchMore={fetchMore}
                  hasMore={hasMore}
                >
                  {isLoading ? (
                    <>
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                      <SkeletonExerciseCard />
                    </>
                  ) : (
                    exercises
                      .filter((exercise) => {
                        if (serachQuery === "") {
                          return true;
                        }

                        return exercise.name
                          .toLowerCase()
                          .includes(serachQuery.toLowerCase());
                      })
                      .map((exercise) => (
                        <ExerciseCard
                          key={exercise.id}
                          exercise={exercise}
                          isSelected={selectedExercisesIds.includes(
                            exercise.id,
                          )}
                          onSelectedChange={() =>
                            toggleExerciseSelection(exercise.id)
                          }
                        />
                      ))
                  )}
                </InfiniteScroll>
              </div>
            </ScrollShadow>
            <div className="h-fit p-2">
              <Button
                className="sticky bottom-0 z-0 w-full"
                color="primary"
                onPress={() => {
                  onSubmit(selectedExercisesIds);
                  onClose();
                }}
              >
                Добавить
              </Button>
            </div>
          </div>
        )}
      </ModalContent>
    </Modal>
  );
}
