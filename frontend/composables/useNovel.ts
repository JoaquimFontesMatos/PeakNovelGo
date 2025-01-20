import { ref } from "vue";
import type { Novel } from "~/models/Novel";

const runtimeConfig = useRuntimeConfig();
const apiUrl = runtimeConfig.public.apiUrl;

export function useNovel(novelTitle: string) {
  const novelData = ref<Novel | null>(null);
  const novelError = ref<string | null>(null);

  async function fetchNovel() {
    try {
      const response = await fetch(
        `${apiUrl}/novels/title/${encodeURIComponent(novelTitle)}`
      );
      if (!response.ok) throw new Error("Failed to fetch novel data");
      novelData.value = await response.json();
    } catch (error: any) {
      novelError.value = error.message || "An error occurred while fetching.";
    }
  }

  return {
    novelData,
    novelError,
    fetchNovel,
  };
}
