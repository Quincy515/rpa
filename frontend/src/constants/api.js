import { SERVER_URL } from "../../config";

const v1 = "api/v1";

export const API = {
  LOGIN: `${SERVER_URL}/${v1}/mini-login`,
  UPDATEUSER: `${SERVER_URL}/${v1}/u/mini-user`,
  UPLOADURL: `${SERVER_URL}/${v1}/upload`
};
