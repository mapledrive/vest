export const addRepos = repos => ({ type: "ADD_REPOS", repos });
export const clearRepos = () => ({ type: "CLEAR_REPOS" });

export const getRepos = username => async dispatch => {
  try {
	var rangestart;
	var rangeend = new Date().toISOString().substring(0, 10);
	if (username === "1")  {rangestart = '2018-01-01'};
	if (username === "2") {rangestart = '2017-01-01'};
	if (username === "5") {rangestart = '2014-01-01'};
	if (username === "10") {rangestart = '2009-01-01'};
	var url = `https://api.thevest.org/vest/v1/charts?code=gc&rs=${rangestart}&re=${rangeend}`;
    const response = await fetch(url);
    const responseBody = await response.json();
    dispatch(addRepos(responseBody));
  } catch (error) {
    console.log(error);
    dispatch(clearRepos());
  }
};


export const addStats = repos => ({ type: "ADD_STATS", repos });
export const getStats = dates => async dispatch => {
  try {
	var rangestart;
	var rangeend = new Date().toISOString().substring(0, 10);
	if (dates === "1")  {rangestart = '2018-01-01'};
	if (dates === "2") {rangestart = '2017-01-01'};
	if (dates === "5") {rangestart = '2014-01-01'};
	if (dates === "10") {rangestart = '2009-01-01'};
	var url = `https://api.thevest.org/vest/v1/stats?code=gc&rangeStart=${rangestart}&rangeEnd=${rangeend}&periodStart=2017-01-01&periodEnd=2018-03-01`;
	const response = await fetch(url);
    const responseBody = await response.json();
    dispatch(addStats(responseBody));
  } catch (error) {
    console.log(error);
    dispatch(clearRepos());
  }
};
