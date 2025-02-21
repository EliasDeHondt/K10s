export const getCookie = (cookieName: string) => new RegExp(`${cookieName}=([^;]+);`).exec(document.cookie)?.[1];