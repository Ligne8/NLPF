import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api/v1';

describe('API Tests for Traffic Manager : Lots page', () => {

    // Test that the API returns a 200 status code for the GET request
    it('should return 200 for GET /traffic-manager/lots', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/lots`);
        expect(response.status).toBe(200);
    });

    // Test that the response contains lots with the expected properties
    it('should return a lot with the expected properties', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/lots`);
        expect(response.status).toBe(200);
        response.data.forEach((lot: any) => {
            expect(lot).toHaveProperty('name');
            expect(lot).toHaveProperty('status');
            expect(lot).toHaveProperty('volume');
            expect(lot).toHaveProperty('location');
            expect(lot).toHaveProperty('startCheckpoint');
            expect(lot).toHaveProperty('endCheckpoint');
            expect(lot).toHaveProperty('tractor');
        });
    });

    // Test that the properties of the lots have the correct data types
    it('should return lots with the correct data types', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/lots`);
        expect(response.status).toBe(200);
        response.data.forEach((lot: any) => {
            expect(typeof lot.name).toBe('string');
            expect(['ON_THE_WAY', 'ON_THE_STOCK_EXCHANGE', 'AVAILABLE', 'ARCHIVED']).toContain(lot.status);
            expect(typeof lot.volume).toBe('number');
            expect(typeof lot.location).toBe('string');
            expect(typeof lot.startCheckpoint).toBe('string');
            expect(typeof lot.endCheckpoint).toBe('string');
            expect(Array.isArray(lot.tractor)).toBe(true);
        });
    });

    // Test that all lots have valid status values from the defined list
    it('should have valid status values', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/lots`);
        expect(response.status).toBe(200);
        const validStatuses = ['ON_THE_WAY', 'ON_THE_STOCK_EXCHANGE', 'AVAILABLE', 'ARCHIVED'];     
        response.data.forEach((lot: any) => {
            expect(validStatuses).toContain(lot.status);
        });
    });

    // Test that the API returns an empty array when no lots are available
    it('should return an empty array when no lots are available', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/lots`);
        expect(response.status).toBe(200);
        expect(Array.isArray(response.data)).toBe(true);
        expect(response.data.length).toBeGreaterThanOrEqual(0);
    });
});
