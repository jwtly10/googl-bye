import axios from 'axios';

const API_BASE_URL = '/v1/api';

const handleResponse = (response) => {
    console.log(response.data);
    return response.data;
};

const handleError = (error) => {
    console.error('API call failed:', error);
    throw error;
};

export const searchGithub = async (searchParams) => {
    try {
        const response = await axios.post(`${API_BASE_URL}/search`, JSON.stringify(searchParams));
        return handleResponse(response);
    } catch (error) {
        return handleError(error);
    }
};
