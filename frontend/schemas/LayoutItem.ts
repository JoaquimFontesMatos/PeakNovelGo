import {z} from 'zod';

const LayoutItemSchema = z.object({
    id: z.string().nonempty('Id is required'),
    type: z.string().nonempty('type is required'),
    props: z.record(z.any()).optional(),
});

type LayoutItem = z.infer<typeof LayoutItemSchema>;
export {type LayoutItem, LayoutItemSchema};