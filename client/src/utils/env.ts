import { z } from "zod";

const envSchema = z.object({
  BASE_URL: z.string().url().nonempty(),
});

export const ENV = envSchema.parse(process.env);
