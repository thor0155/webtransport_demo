import { config } from "@/config";
import axios from "axios";

export const api = axios.create({
    baseURL: config.apiBaseUrl,
    timeout: 10000,
});

export async function GetCertHash() {
    const response = await api.get<string>("/cert", {
        responseType: "text",
    });

    if (response.status !== 200) {
        throw new Error(`HTTP ${response.status}`);
    }

    return response.data;
}
