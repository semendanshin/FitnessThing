"use client";
import React from "react";
import { Form, Input, Button } from "@nextui-org/react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";

import { authApi } from "@/api/api";

export default function LoginPage() {
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");

  const [errors, setErrors] = React.useState({});

  const router = useRouter();

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (email === "") {
      setErrors({ email: "Email необходим" });

      return;
    }

    if (password === "") {
      setErrors({ password: "Пароль необходим" });

      return;
    }

    async function fetchData() {
      authApi.v1
        .authServiceLogin({
          email: email.toString(),
          password: password.toString(),
        })
        .then((response) => {
          console.log(response.data);
          if (response.data.tokens !== undefined) {
            localStorage.setItem(
              "accessToken",
              response.data.tokens.accessToken,
            );
            localStorage.setItem(
              "refreshToken",
              response.data.tokens.refreshToken,
            );
            router.push("/");
          }
        })
        .catch((error) => {
          console.log(error);
          if (error.response?.status === 400) {
            setErrors({ password: "Неверный email или пароль" });
          }
          toast.error("Ошибка входа");
        });
    }

    fetchData();
    setErrors({});
  };

  return (
    <>
      <div className="align-left w-full">
        <h1 className="text-4xl font-bold leading-tight">Вход</h1>
        <p className="text-gray-600">Войдите в свой аккаунт</p>
      </div>
      <Form
        className="inline-block text-center justify-center w-full max-w-lg"
        validationBehavior="native"
        validationErrors={errors}
        onSubmit={handleSubmit}
      >
        <div className="flex flex-col items-center justify-center gap-4 py-4">
          <Input
            isRequired
            autoComplete="email"
            label="Email"
            labelPlacement="outside"
            name="email"
            placeholder="Email"
            type="email"
            value={email}
            onValueChange={(value) => setEmail(value)}
          />
          <Input
            isRequired
            autoComplete="current-password"
            label="Пароль"
            labelPlacement="outside"
            name="password"
            placeholder="Пароль"
            type="password"
            value={password}
            onValueChange={(value) => setPassword(value)}
          />
          <div className="flex items-center gap-1">
            <p className="text-gray-600">Нет аккаунта?</p>
            <p className="text-primary">
              {<Link href="/auth/register">Регистрация</Link>}
            </p>
          </div>
          <Button color="primary" type="submit">
            Войти
          </Button>
        </div>
      </Form>
    </>
  );
}
