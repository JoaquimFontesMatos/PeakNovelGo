import { z } from 'zod';

const ToastTypeSchema = z.enum(['error', 'success', 'warning']);
const ToastIconSchema = z.enum(['user', 'project', 'auth', 'none', 'chapter', 'novel', 'bookmark', 'tts']);

const ToastSchema = z.object({
  id: z.number(),
  message: z.string(),
  type: ToastTypeSchema,
  icon: ToastIconSchema.optional(),
});

type Toast = z.infer<typeof ToastSchema>;
type ToastType = z.infer<typeof ToastTypeSchema>;
type ToastIcon = z.infer<typeof ToastIconSchema>;

export { type Toast, type ToastType, type ToastIcon, ToastSchema, ToastTypeSchema, ToastIconSchema };
