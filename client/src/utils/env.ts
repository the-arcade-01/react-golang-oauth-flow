import { z } from "zod";

const envSchema = z.object({
  VITE_API_URL: z.string(),
});

const ENV = envSchema.safeParse(import.meta.env);

if (!ENV.success) {
  console.error("Invalid environment variables:", ENV.error.format());
  throw new Error("Invalid environment variables");
}

export const env = ENV.data;
