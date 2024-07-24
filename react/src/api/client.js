import axios from 'axios';

const API_BASE_URL = '/v1/api';

const handleResponse = (response) => {
    console.log('Request Response: ', response.data);
    return response.data;
};

const handleError = (error) => {
    console.error('API call failed:', error);
    throw error;
};

export const searchGithub = async (searchParams) => {
    try {
        // Create a new object with the same properties as searchParams
        const processedParams = { ...searchParams };

        // Convert specific fields to integers
        processedParams.startPage = parseInt(processedParams.startPage, 10);
        processedParams.currentPage = parseInt(processedParams.currentPage, 10);
        processedParams.pagesToProcess = parseInt(processedParams.pagesToProcess, 10);

        const response = await axios.post(`${API_BASE_URL}/search`, JSON.stringify(processedParams));
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};

export const saveRepos = async (reposFromGithub) => {
    // We need to do some parsing as the repos from the golang server have the embedded model fields.
    // We need to remove them before trying to post to the golang API

    //     var body = [];
    //     reposFromGithub.forEach((repo, _) => {
    //         const item = {
    //             name: repo.name,
    //             author: repo.author,
    //             language: repo.language,
    //             stars: repo.stars,
    //             forks: repo.forks,
    //             size: repo.size,
    //             lastPush: repo.lastPush,
    //             apiUrl: repo.apiUrl,
    //             ghUrl: repo.ghUrl,
    //             cloneUrl: repo.cloneUrl,
    //         };
    //         body.push(item);
    //     });

    try {
        const response = await axios.post(`${API_BASE_URL}/save`, JSON.stringify(reposFromGithub));
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};
