import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api/v1';

describe('API Tests', () => {
    it('should return 200 for GET /checkpoints', async () => {
        const response = await axios.get(`${BASE_URL}/checkpoints`);
        expect(response.status).toBe(200);
    });
});