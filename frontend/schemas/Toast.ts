import { z } from 'zod';

const ToastTypeSchema = z.enum(['error', 'success', 'warning']);
const ToastIconSchema = z.enum(['user', 'project', 'auth', 'none', 'chapter', 'novel', 'bookmark', 'tts', 'network']);
const ToastIconMap: Record<ToastIcon, string> = {
  network: 'fluent:network-check-24-filled',
  user: 'fluent:person-24-filled',
  project: 'fluent:pulse-24-filled',
  auth: 'fluent:shield-lock-24-filled',
  none: 'none',
  chapter: 'fluent:document-one-page-24-filled',
  novel: 'fluent:book-24-filled',
  bookmark: 'fluent:bookmark-multiple-24-filled',
  tts: 'fluent:desktop-speaker-24-filled',
};

const ToastSchema = z.object({
  id: z.number(),
  message: z.string(),
  type: ToastTypeSchema,
  icon: z.string(),
});

type Toast = z.infer<typeof ToastSchema>;
type ToastType = z.infer<typeof ToastTypeSchema>;
type ToastIcon = z.infer<typeof ToastIconSchema>;

export { type Toast, type ToastType, type ToastIcon, ToastIconMap, ToastSchema, ToastTypeSchema, ToastIconSchema };
