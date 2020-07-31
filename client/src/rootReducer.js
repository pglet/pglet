import { combineReducers } from "redux";
import pageReducer from './features/page/pageSlice'
import userReducer from './features/users/userSlice'

export default combineReducers({
  page: pageReducer,
  user: userReducer
});
