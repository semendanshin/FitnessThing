"use client";
import React from "react";
import { Form, Input, Button } from "@nextui-org/react";
import Link from "next/link";
import { authApi } from "@/api/api";
import { useRouter } from "next/navigation";

export default function LoginPage() {
    const [email, setEmail] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [firstName, setFirstName] = React.useState("");
    const [lastName, setLastName] = React.useState("");
    const [dateOfBirth, setDateOfBirth] = React.useState("");

    const [errors, setErrors] = React.useState<{
        email?: string;
        password?: string;
        passwordRepeated?: string;
        general?: string;
    }>({});

    const router = useRouter();

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const data = Object.fromEntries(new FormData(e.currentTarget));

        if (data.email === "") {
            setErrors({ email: "Email необходим" });
            return;
        }

        if (data.password === "") {
            setErrors({ password: "Пароль необходим" });
            return;
        }

        if (data.password !== data.passwordRepeated) {
            setErrors({ passwordRepeated: "Пароли не совпадают" });
            return;
        }

        setErrors({});

        async function fetchData() {
            authApi.v1
                .userServiceCreateUser({
                    email: data.email.toString(),
                    password: data.password.toString(),
                    firstName: data.firstName.toString(),
                    lastName: data.lastName.toString(),
                })
                .then((response) => {
                    console.log(response.data);
                    router.push("/auth/login");
                })
                .catch((error) => {
                    console.log(error);

                    if (error.response?.status === 400) {
                        setErrors({
                            general:
                                "Ошибка в запросе. Обратитесь к администратору",
                        });
                    }

                    if (error.response?.status === 409) {
                        setErrors({ email: "Такой email уже зарегистрирован" });
                    }

                    if (error.response?.status === 422) {
                        setErrors({ general: "Ошибка валидации" });
                    }

                    if (error.response?.status === 500) {
                        setErrors({ general: "Ошибка сервера" });
                    }

                    return;
                });
        }

        fetchData();
    };

    return (
        <>
            <div className="align-left w-full">
                <h1 className="text-4xl font-bold leading-tight">
                    Регистрацияя
                </h1>
                <p className="text-gray-600">Создайте новый аккаунт</p>
            </div>
            <Form
                validationBehavior="aria"
                validationErrors={errors}
                className="inline-block text-center justify-center w-full max-w-lg"
                onSubmit={handleSubmit}
            >
                <div className="flex flex-col items-center justify-center gap-4 py-4">
                    <Input
                        isRequired
                        label="Email"
                        placeholder="Email"
                        labelPlacement="outside"
                        name="email"
                        type="email"
                        autoComplete="email"
                    />
                    <Input
                        isRequired
                        label="Пароль"
                        placeholder="Пароль"
                        labelPlacement="outside"
                        name="password"
                        type="password"
                        validate={(value) => {
                            if (value && value.length < 8) {
                                return "Пароль должен быть длиннее 8 символов";
                            }

                            return undefined;
                        }}
                    />
                    <Input
                        isRequired
                        label="Повторите пароль"
                        placeholder="Пароль"
                        labelPlacement="outside"
                        name="passwordRepeated"
                        autoComplete="new-password"
                        type="password"
                    />
                    <Input
                        label="Имя"
                        placeholder="Иван"
                        labelPlacement="outside"
                        name="firstName"
                        type="text"
                    />
                    <Input
                        label="Фамилия"
                        placeholder="Иванов"
                        name="lastName"
                        labelPlacement="outside"
                        type="text"
                    />
                    {errors.general ? (
                        <div className="flex items-center gap-1 w-full">
                            <p className="text-red-500 text-sm">
                                {errors.general}
                            </p>
                        </div>
                    ) : null}
                    <div className="flex items-center gap-1">
                        <p className="text-gray-600">Уже есть аккаунт?</p>
                        <p className="text-primary">
                            {<Link href="/auth/login">Войти</Link>}
                        </p>
                    </div>
                    <Button color="primary" type="submit">
                        Зарегистрироваться
                    </Button>
                </div>
            </Form>
        </>
    );
}
