import * as yup from 'yup';

// Define the Yup schema for the TTS part of ReadingPreferences
const TtsSchema = yup.object({
    autoplay: yup.boolean().nullable().notRequired(),
    voice: yup.string().nullable().notRequired(),
    rate: yup.number().nullable().notRequired(),
});

// Define the Yup schema for ReadingPreferences
const ReadingPreferencesSchema = yup.object({
    atomicReading: yup.boolean().nullable().notRequired(),
    font: yup.string().nullable().notRequired(),
    theme: yup.string().nullable().notRequired(),
    tts: TtsSchema.required('TTS settings are required'),
});

export { TtsSchema, ReadingPreferencesSchema };