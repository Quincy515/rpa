import { LOGIN } from "@/constants/user";

const INITIAL_STATE = {
  data: null
};

export default function user(state = INITIAL_STATE, action) {
  switch (action.type) {
    case LOGIN:
      return {
        ...state,
        data: action.payload
      };
    default:
      return state;
  }
}
