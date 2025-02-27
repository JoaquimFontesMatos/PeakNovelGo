import type { TextUtils } from '~/interfaces/TextUtils';

export class BaseTextUtils implements TextUtils {
  convertLineBreaksToHtml(text: string): string {
    return text.replace(/\n/g, '<br>');
  }

  toBionicText(text: string): string {
    return text
      .split(' ') // Split the text into words
      .map(word => {
        if (word.length <= 2) {
          // Entire word is emphasized if it's very short
          return `<b>${word}</b>`;
        } else {
          const boldLength = Math.ceil(word.length / 2); // Emphasize half of the word
          const boldPart = word.slice(0, boldLength);
          const restPart = word.slice(boldLength);
          return `<b>${boldPart}</b>${restPart}`;
        }
      })
      .join(' '); // Rejoin the words into a single string
  }
}
