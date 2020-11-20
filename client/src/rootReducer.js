import { combineReducers } from "redux";
import pageReducer from './slices/pageSlice'
import userReducer from './slices/userSlice'

export default combineReducers({
  page: pageReducer,
  user: userReducer
});
