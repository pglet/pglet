import { combineReducers } from "redux";
import page from './reducers/page'
import user from './reducers/user'

export default combineReducers({
  page,
  user
});
