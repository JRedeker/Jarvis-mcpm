/**
 * Clients Routes Tests
 */

const request = require('supertest');
const { createApp } = require('../server');

describe('Clients Routes', () => {
    let app;

    beforeAll(() => {
        app = createApp();
    });

    describe('GET /api/v1/clients', () => {
        it('should return list of clients', async () => {
            const res = await request(app)
                .get('/api/v1/clients')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.count).toBeGreaterThan(0);
            expect(Array.isArray(res.body.data.clients)).toBe(true);
        });

        it('should include known clients', async () => {
            const res = await request(app)
                .get('/api/v1/clients')
                .expect(200);

            const clientNames = res.body.data.clients.map(c => c.name);
            expect(clientNames).toContain('opencode');
            expect(clientNames).toContain('claude-code');
            expect(clientNames).toContain('claude-desktop');
        });

        it('should include client details', async () => {
            const res = await request(app)
                .get('/api/v1/clients')
                .expect(200);

            const client = res.body.data.clients[0];
            expect(client.name).toBeDefined();
            expect(client.displayName).toBeDefined();
            expect(typeof client.detected).toBe('boolean');
        });
    });

    describe('GET /api/v1/clients/:name', () => {
        it('should return client info for known client', async () => {
            const res = await request(app)
                .get('/api/v1/clients/opencode')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.name).toBe('opencode');
            expect(res.body.data.displayName).toBe('OpenCode');
            expect(res.body.data.possiblePaths).toBeDefined();
        });

        it('should return 404 for unknown client', async () => {
            const res = await request(app)
                .get('/api/v1/clients/unknown-client-xyz')
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('CLIENT_NOT_FOUND');
            expect(res.body.error.message).toContain('Known clients');
        });
    });

    describe('PUT /api/v1/clients/:name', () => {
        it('should return 404 for unknown client', async () => {
            const res = await request(app)
                .put('/api/v1/clients/unknown-client-xyz')
                .send({ add_profile: 'test' })
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('CLIENT_NOT_FOUND');
        });

        it('should return 400 for no changes', async () => {
            const res = await request(app)
                .put('/api/v1/clients/opencode')
                .send({})
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('NO_CHANGES');
        });
    });
});
