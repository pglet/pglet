import { combineReducers } from "redux";
import pageReducer from './features/page/pageSlice'
import user from './reducers/user'

export default combineReducers({
  page: pageReducer,
  user
});
