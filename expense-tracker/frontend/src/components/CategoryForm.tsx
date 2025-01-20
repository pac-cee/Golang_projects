import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Chip,
  IconButton,
} from '@mui/material';
import { Add as AddIcon, Close as CloseIcon } from '@mui/icons-material';
import { Category } from '../types';
import { useAppContext } from '../context/AppContext';
import { categoryApi } from '../services/api';

interface CategoryFormProps {
  open: boolean;
  onClose: () => void;
  category?: Category;
}

const CategoryForm = ({ open, onClose, category }: CategoryFormProps) => {
  const { dispatch } = useAppContext();
  const [name, setName] = useState('');
  const [newSubcategory, setNewSubcategory] = useState('');
  const [subcategories, setSubcategories] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (category) {
      setName(category.name);
      setSubcategories(category.subcategories);
    }
  }, [category]);

  const handleAddSubcategory = () => {
    if (newSubcategory.trim() && !subcategories.includes(newSubcategory.trim())) {
      setSubcategories([...subcategories, newSubcategory.trim()]);
      setNewSubcategory('');
    }
  };

  const handleDeleteSubcategory = (subcategory: string) => {
    setSubcategories(subcategories.filter((sub) => sub !== subcategory));
  };

  const handleSubmit = async () => {
    if (!name.trim()) return;

    try {
      setLoading(true);
      const data = {
        name: name.trim(),
        subcategories,
      };

      if (category) {
        const response = await categoryApi.update(category.id, data);
        dispatch({ type: 'UPDATE_CATEGORY', payload: response.data });
      } else {
        const response = await categoryApi.create(data);
        dispatch({ type: 'ADD_CATEGORY', payload: response.data });
      }
      onClose();
    } catch (error) {
      console.error('Error saving category:', error);
      dispatch({ type: 'SET_ERROR', payload: 'Error saving category' });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setName('');
    setNewSubcategory('');
    setSubcategories([]);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {category ? 'Edit Category' : 'Add New Category'}
      </DialogTitle>
      <DialogContent>
        <Box sx={{ mt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            fullWidth
            label="Category Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <Box sx={{ display: 'flex', gap: 1 }}>
            <TextField
              fullWidth
              label="Add Subcategory"
              value={newSubcategory}
              onChange={(e) => setNewSubcategory(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === 'Enter') {
                  handleAddSubcategory();
                }
              }}
            />
            <IconButton
              color="primary"
              onClick={handleAddSubcategory}
              disabled={!newSubcategory.trim()}
            >
              <AddIcon />
            </IconButton>
          </Box>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {subcategories.map((subcategory) => (
              <Chip
                key={subcategory}
                label={subcategory}
                onDelete={() => handleDeleteSubcategory(subcategory)}
                deleteIcon={<CloseIcon />}
              />
            ))}
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={loading || !name.trim()}
        >
          {loading ? 'Saving...' : 'Save'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CategoryForm;
