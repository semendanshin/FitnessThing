"use client";

import type { Edge } from "@atlaskit/pragmatic-drag-and-drop-hitbox/types";

import { Button } from "@nextui-org/button";
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { DropdownItem } from "@nextui-org/dropdown";
import { Form } from "@nextui-org/form";
import { Input, Textarea } from "@nextui-org/input";
import {
  Modal,
  ModalContent,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/modal";
import { useRouter } from "next/navigation";
import { use, useEffect, useRef, useState } from "react";
import { toast } from "react-toastify";
import {
  draggable,
  dropTargetForElements,
} from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { DropIndicator } from "@atlaskit/pragmatic-drag-and-drop-react-drop-indicator/box";
import { combine } from "@atlaskit/pragmatic-drag-and-drop/combine";
import {
  attachClosestEdge,
  extractClosestEdge,
} from "@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge";
import { reorder } from "@atlaskit/pragmatic-drag-and-drop/reorder";
import invariant from "tiny-invariant";
import clsx from "clsx";
import Link from "next/link";

import { ChevronRightIcon, GripVerticalIcon, PlusIcon } from "@/config/icons";
import { ModalSelectExercise } from "@/components/pick-exercises-modal";
import { PageHeader } from "@/components/page-header";
import { Loading } from "@/components/loading";
import {
  WorkoutExerciseInstanceDetails,
  WorkoutRoutineDetailResponse,
} from "@/api/api.generated";
import { authApi, errUnauthorized } from "@/api/api";
// import "@/scripts/drag-drop-touch-patch";
// import "drag-drop-touch";

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const [isLoading, setIsLoading] = useState(true);

  const [routineDetails, setRoutineDetails] =
    useState<WorkoutRoutineDetailResponse>({});

  const { id } = use(params);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const {
    isOpen: isRenameOpen,
    onOpen: onRenameOpen,
    onOpenChange: onRenameOpenChange,
  } = useDisclosure();

  const router = useRouter();

  async function fetchRoutineDetails() {
    console.log(id);
    await authApi.v1
      .routineServiceGetRoutineDetail(id)
      .then((response) => {
        console.log(response.data);
        setRoutineDetails(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (
          error === errUnauthorized ||
          (error as any).response?.status === 401
        ) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при загрузке данных");
      });
  }

  async function addExerciseToRoutine(exerciseId: string) {
    await authApi.v1
      .routineServiceAddExerciseToRoutine(id, { exerciseId: exerciseId })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при добавлении упражнения");
      });
  }

  async function submitPickExercise(exerciseIds: string[]) {
    console.log(exerciseIds);
    try {
      for (const exerciseId of exerciseIds) {
        await addExerciseToRoutine(exerciseId);
      }
      await fetchRoutineDetails();
    } catch (error) {
      console.log(error);
      if (
        error === errUnauthorized ||
        (error as any).response?.status === 401
      ) {
        router.push("/auth/login");

        return;
      }
      toast.error("Ошибка при добавлении упражнения");
    }
  }

  async function reoderExerciseInstances(
    sourceIndex: number,
    targetIndex: number,
  ) {
    let newExerciseInstances = reorder({
      list: routineDetails.exerciseInstances!,
      startIndex: sourceIndex,
      finishIndex: targetIndex,
    });

    setRoutineDetails((prev) => ({
      ...prev,
      exerciseInstances: newExerciseInstances,
    }));

    const exerciseInstanceIds = newExerciseInstances.map(
      (ei) => ei.exerciseInstance!.id!,
    );

    await authApi.v1
      .routineServiceSetExerciseOrder(id, {
        exerciseInstanceIds: exerciseInstanceIds!,
      })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при изменении порядка упражнений");
      });
  }

  async function fetchData() {
    setIsLoading(true);
    await fetchRoutineDetails();
    setIsLoading(false);
  }

  useEffect(() => {
    fetchData();
  }, []);

  if (isLoading) {
    return <Loading />;
  }

  function ExercieInstanceCard({
    exerciseInstanceDetails,
  }: {
    exerciseInstanceDetails: WorkoutExerciseInstanceDetails;
  }) {
    const setsCount = exerciseInstanceDetails.sets?.length || 0;

    const ref = useRef<HTMLDivElement>(null);
    const dragHandleRef = useRef(null);
    const [isDragged, setIsDragged] = useState<boolean>(false);
    const [closestEdge, setClosestEdge] = useState<Edge | null>(null);

    const [data] = useState({
      id: exerciseInstanceDetails.exerciseInstance!.id,
    });

    useEffect(() => {
      const el = ref.current;
      const dragHandleEl = dragHandleRef.current;

      invariant(el);
      invariant(dragHandleEl);

      return combine(
        draggable({
          element: el,
          getInitialData: () => data,
          onDragStart: () => {
            setIsDragged(true);

            if (navigator.vibrate) {
              navigator.vibrate(100);
            }
          },
          onDrop: () => setIsDragged(false),
        }),
        dropTargetForElements({
          element: el,
          canDrop: ({ source }) => {
            return (
              source.data.id !== exerciseInstanceDetails.exerciseInstance!.id
            );
          },
          getData: ({ input }) => {
            return attachClosestEdge(data, {
              element: el,
              input,
              allowedEdges: ["top", "bottom"],
            });
          },
          onDrag: ({ self, source }) => {
            const isSource = source.element === el;

            if (isSource) {
              setClosestEdge(null);

              return;
            }

            const closestEdgeValue = extractClosestEdge(self.data);

            const exerciseInstanceIds = routineDetails.exerciseInstances?.map(
              (ei) => ei.exerciseInstance!.id,
            );

            const sourceIndex = exerciseInstanceIds?.indexOf(
              source.data.id as string,
            )!;

            const currentIndex = exerciseInstanceIds?.indexOf(
              exerciseInstanceDetails.exerciseInstance!.id,
            );

            const isItemBeforeSource = currentIndex === sourceIndex - 1;
            const isItemAfterSource = currentIndex === sourceIndex + 1;

            const isDropIndicatorHidden =
              (isItemBeforeSource && closestEdgeValue === "bottom") ||
              (isItemAfterSource && closestEdgeValue === "top");

            if (isDropIndicatorHidden) {
              setClosestEdge(null);

              return;
            }

            setClosestEdge(closestEdgeValue);
          },
          onDragLeave: () => {
            setClosestEdge(null);
          },
          onDrop: ({ source, self }) => {
            setClosestEdge(null);

            const exerciseInstanceIds = routineDetails.exerciseInstances?.map(
              (ei) => ei.exerciseInstance!.id,
            );

            const sourceIndex = exerciseInstanceIds?.indexOf(
              source.data.id as string,
            )!;

            let currentIndex = exerciseInstanceIds?.indexOf(
              exerciseInstanceDetails.exerciseInstance!.id,
            )!;

            const closestEdgeValue = extractClosestEdge(self.data);

            const isItemBeforeSource = currentIndex < sourceIndex;
            const isItemAfterSource = currentIndex > sourceIndex;

            if (closestEdgeValue === "bottom" && isItemBeforeSource) {
              currentIndex += 1;
            } else if (closestEdgeValue === "top" && isItemAfterSource) {
              currentIndex -= 1;
            }

            console.log(sourceIndex, currentIndex);

            reoderExerciseInstances(sourceIndex, currentIndex);
          },
        }),
      );
    }, []);

    return (
      <div className="relative">
        <Card
          key={exerciseInstanceDetails.exerciseInstance!.id}
          ref={ref}
          fullWidth
          className={clsx(
            "flex flex-row flex-grow p-3 gap-2 justify-between select-none",
            isDragged && "transform scale-95",
          )}
          isDisabled={isDragged}
          shadow="sm"
          onContextMenu={(e) => {
            e.preventDefault();
          }}
          onPress={() => {
            router.push(
              `/routines/${id}/exerciseInstance/${exerciseInstanceDetails.exerciseInstance!.id}`,
            );
          }}
        >
          <div className="flex flex-row gap-3">
            <div className="flex flex-col items-start justify-center">
              <div className="rounded-md bg-default-100 hover:bg-default-200 cursor-grab flex items-center justify-center p-1">
                <GripVerticalIcon ref={dragHandleRef} className="w-4 h-4" />
              </div>
            </div>
            <div className="flex flex-col items-start justify-between">
              <CardHeader className="p-0">
                <p className="text-m font-bold text-start">
                  {exerciseInstanceDetails.exercise!.name}
                </p>
              </CardHeader>
              <CardBody className="p-0">
                <div className="flex flex-row gap-1">
                  <p className="text-xs text-gray-400 whitespace-nowrap">
                    {setsCount}{" "}
                    {"подход" +
                      (setsCount % 10 === 1
                        ? ""
                        : setsCount % 10 >= 2 && setsCount % 10 <= 4
                          ? "а"
                          : "ов")}
                    {" •"}
                  </p>
                  <p className="text-xs text-gray-400 whitespace-nowrap">
                    {exerciseInstanceDetails.exercise!.targetMuscleGroups!.join(
                      ", ",
                    )}
                  </p>
                </div>
              </CardBody>
            </div>
          </div>
          <Link
            className="flex flex-col items-center justify-center"
            href={`/routines/${id}/exerciseInstance/${exerciseInstanceDetails.exerciseInstance!.id}`}
          >
            <ChevronRightIcon className="w-4 h-4" fill="currentColor" />
          </Link>
        </Card>
        {closestEdge && <DropIndicator edge={closestEdge} gap="1rem" />}
      </div>
    );
  }

  function ModalRenameRoutine() {
    const [routineName, setRoutineName] = useState(
      routineDetails.routine?.name,
    );
    const [routineDescription, setRoutineDescription] = useState(
      routineDetails.routine?.description,
    );
    const [isButtonLoading, setIsButtonLoading] = useState(false);

    async function updateRoutine() {
      setIsButtonLoading(true);
      try {
        await authApi.v1.routineServiceUpdateRoutine(id, {
          name: routineName,
          description: routineDescription,
        });
        await fetchRoutineDetails();
      } catch (error) {
        console.log(error);
        toast.error("Ошибка при обновлении рутины");
        if (
          error === errUnauthorized ||
          (error as any).response?.status === 401
        ) {
          router.push("/auth/login");

          return;
        }
      } finally {
        setIsButtonLoading(false);
      }
    }

    return (
      <Modal
        className="overflow-y-auto p-2 w-full"
        isOpen={isRenameOpen}
        placement="center"
        scrollBehavior="inside"
        size="xs"
        onClose={onRenameOpenChange}
      >
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="p-2">Переименовать рутину</ModalHeader>

              <Form
                className="inline-block text-center justify-center w-full max-w-lg p-0"
                onSubmit={async (e) => {
                  e.preventDefault();
                  await updateRoutine();
                  onClose();
                }}
              >
                <div className="grid grid-cols-1 gap-4 p-2">
                  <Input
                    autoFocus
                    placeholder="Название"
                    value={routineName}
                    onChange={(e) => setRoutineName(e.target.value)}
                  />
                  <Textarea
                    placeholder="Описание"
                    value={routineDescription ? routineDescription : ""}
                    onChange={(e) => setRoutineDescription(e.target.value)}
                  />
                  <Button
                    color="primary"
                    isLoading={isButtonLoading}
                    type="submit"
                  >
                    Сохранить
                  </Button>
                </div>
              </Form>
            </>
          )}
        </ModalContent>
      </Modal>
    );
  }

  async function onDelete() {
    await authApi.v1
      .routineServiceDeleteRoutine(id)
      .catch((error) => {
        console.log(error);
        if (
          error === errUnauthorized ||
          (error as any).response?.status === 401
        ) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при удалении рутины");
      })
      .then(() => {
        router.back();
      });
  }

  return (
    <>
      <div className="flex flex-col py-4 gap-4">
        <PageHeader enableBackButton title={routineDetails.routine?.name!}>
          <DropdownItem key="rename" color="primary" onPress={onRenameOpen}>
            Переименовать
          </DropdownItem>
          <DropdownItem
            key="delete"
            className="text-danger"
            color="danger"
            onPress={onDelete}
          >
            Удалить
          </DropdownItem>
        </PageHeader>
        {routineDetails.routine?.description ? (
          <div className="px-4">
            <p className="text-sm text-gray-500">
              {routineDetails.routine?.description}
            </p>
          </div>
        ) : null}
        <div className="grid grid-cols-1 px-4 gap-4">
          {routineDetails.exerciseInstances?.map(
            (exerciseInstanceDetails: WorkoutExerciseInstanceDetails) => (
              <ExercieInstanceCard
                key={exerciseInstanceDetails.exerciseInstance!.id}
                exerciseInstanceDetails={exerciseInstanceDetails}
              />
            ),
          )}
          <Card className="p-2">
            <Button onPress={onOpen}>
              <PlusIcon className="w-4 h-4" />
              Добавить упражнение
            </Button>
          </Card>
        </div>
      </div>
      <ModalSelectExercise
        excludeExerciseIds={routineDetails.exerciseInstances!.map(
          (ei) => ei.exercise!.id!,
        )}
        isOpen={isOpen}
        onClose={onClose}
        onSubmit={submitPickExercise}
      />

      <ModalRenameRoutine />
    </>
  );
}
