"use client";
import { authApi, errUnauthorized } from "@/api/api";
import { Loading } from "@/components/loading";
import {
    Avatar,
    Button,
    Card,
    CardBody,
    Divider,
    Spinner,
} from "@nextui-org/react";
import { useTheme } from "next-themes";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

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
                if (
                    error === errUnauthorized ||
                    error.response?.status === 401
                ) {
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
                <h2 className="text-lg text-red-500">
                    Ошибка при загрузке данных
                </h2>
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

    return (
        <div className="p-4 flex-grow gap-4">
            {/* <h1 className="text-2xl font-bold leading-tight">Профиль</h1> */}
            <div className="grid grid-cols-1 gap-4 py-4">
                <div className="flex flex-col gap-4 items-center justify-around">
                    <Avatar size="lg" src="https://i.pravatar.cc/300" />
                    <h2 className="text-2xl font-bold">
                        {user.firstName} {user.lastName}
                    </h2>
                </div>
                <Card className="flex flex-row flex-grow p-0 gap-4 justify-between">
                    <CardBody>
                        <p>
                            <b>Email:</b>
                        </p>
                        <p>{user.email}</p>
                    </CardBody>
                </Card>
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
                <Button onPress={handleLogout} color="danger">
                    Выйти
                </Button>
            </div>
        </div>
    );
}
