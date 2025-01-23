import { Api, ApiConfig, WorkoutTokensPair } from "./api.generated";

const ACCESS_TOKEN_KEY = "accessToken";
const REFRESH_TOKEN_KEY = "refreshToken";
const HEADER_AUTHORIZATION = "X-Access-Token";

export const errUnauthorized = new Error("Unauthorized");

class AuthApi<T = any> extends Api<T> {
    static isRefreshing = false;
    static refreshSubscribers: Array<(token: string) => void> = [];

    config: ApiConfig;

    constructor(config: ApiConfig = {}) {
        super(config);
        this.config = config;
        // Добавляем interceptor для запросов
        this.instance.interceptors.request.use(
            (config) => {
                const token = localStorage.getItem(ACCESS_TOKEN_KEY);

                if (token) {
                    config.headers[HEADER_AUTHORIZATION] = token;
                }

                return config;
            },
            (error) => Promise.reject(error)
        );

        // Добавляем interceptor для ответов
        this.instance.interceptors.response.use(
            (response) => response,
            async (error) => {
                const originalRequest = error.config;

                // Если ошибка 401 и это не запрос на обновление токена
                if (error.response?.status === 401 && !originalRequest._retry) {
                    console.log("Unauthorized, trying to refresh token");
                    originalRequest._retry = true;

                    if (!AuthApi.isRefreshing) {
                        AuthApi.isRefreshing = true;
                        try {
                            const refreshToken =
                                localStorage.getItem(REFRESH_TOKEN_KEY);

                            if (!refreshToken) {
                                throw new Error("No refresh token");
                            }

                            const accessToken =
                                localStorage.getItem(ACCESS_TOKEN_KEY);

                            if (!accessToken) {
                                throw new Error("No access token");
                            }

                            // Создаем временный клиент для обновления токена
                            const tempApi = new Api(config);
                            const response =
                                await tempApi.v1.authServiceRefresh({
                                    tokens: {
                                        accessToken,
                                        refreshToken,
                                    },
                                });

                            const tokens = response.data
                                .tokens as WorkoutTokensPair;

                            this.updateTokens(tokens);

                            console.log("Token refreshed, retrying request");
                            originalRequest.headers[HEADER_AUTHORIZATION] =
                                tokens.accessToken;

                            AuthApi.refreshSubscribers.forEach((cb) =>
                                cb(tokens.accessToken)
                            );
                            AuthApi.refreshSubscribers = [];

                            return this.instance(originalRequest);
                        } catch (refreshError) {
                            console.log(
                                "Failed to refresh token",
                                refreshError
                            );
                            localStorage.removeItem(ACCESS_TOKEN_KEY);
                            localStorage.removeItem(REFRESH_TOKEN_KEY);

                            return Promise.reject(errUnauthorized);
                        } finally {
                            AuthApi.isRefreshing = false;
                        }
                    } else {
                        return new Promise((resolve) => {
                            AuthApi.refreshSubscribers.push((token: string) => {
                                originalRequest.headers[HEADER_AUTHORIZATION] =
                                    token;
                                resolve(this.instance(originalRequest));
                            });
                        });
                    }
                }

                return Promise.reject(error);
            }
        );
    }

    private updateTokens(tokens: WorkoutTokensPair) {
        localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken);
        localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken);
    }
}

export const authApi = new AuthApi({
    baseURL: process.env.NEXT_PUBLIC_API_URL || "/api",
});
