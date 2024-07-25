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

export const searchGithubRepos = async (searchParams) => {
    try {
        // Create a new object with the same properties as searchParams
        const processedParams = { ...searchParams };

        // Convert specific fields to integers
        processedParams.startPage = parseInt(processedParams.startPage, 10);
        processedParams.currentPage = parseInt(processedParams.currentPage, 10);
        processedParams.pagesToProcess = parseInt(processedParams.pagesToProcess, 10);

        const response = await axios.post(
            `${API_BASE_URL}/search`,
            JSON.stringify(processedParams)
        );
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};

export const searchGithubUsers = async (username) => {
    try {
        const response = await axios.get(`${API_BASE_URL}/search-user`, {
            params: { username },
        });
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};

export const searchGithubUsersRepo = async (username) => {
    try {
        const response = await axios.get(`${API_BASE_URL}/search`, {
            params: { username },
        });
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};

export const saveRepos = async (reposFromGithub) => {
    try {
        const response = await axios.post(`${API_BASE_URL}/save`, JSON.stringify(reposFromGithub));
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};

export const searchRepoLinks = async () => {
    try {
        const response = await axios.get(`${API_BASE_URL}/repoLinks`);
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};

export const searchRepoLinksForUser = async (username) => {
    try {
        const response = await axios.get(`${API_BASE_URL}/repoLinks-user`, {
            params: { username },
        });
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};
